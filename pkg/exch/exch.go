package exch

import (
	"errors"
)

type smChan chan int
type ifChan chan interface{}
type ifRead func() (interface{}, error)
type ifWrite func(interface{}) error

// Exchange a pair of read and write with io chan
func Exchange(
	source chan interface{},
	sink chan interface{},
	read func() (interface{}, error),
	write func(interface{}) error,
	sema chan int,
	coInit func(bool) error,
) error {
	// Check semaphore, if taken we quit
	if sema != nil {
		select {
		case sema <- 0:
			defer func() {
				<-sema
			}()
			break
		default:
			return errors.New("channel is not open")
		}
	}

	// Run co-init
	if coInit != nil {
		err := coInit(true)
		if err != nil {
			return err
		}
		defer coInit(false)
	}

	// Start the source write loop
	stop := make(smChan)
	defer func() {
		stop <- 0
	}()
	go sourceWriteLoop(stop, source, write)

	// Read -> Sink loop
	for {
		data, err := read()

		// If error, we quit
		if err != nil {
			return nil
		}

		// Write to sink if we can
		select {
		case sink <- data:
		default:
		}
	}
}

func sourceWriteLoop(
	stop smChan,
	source ifChan,
	write ifWrite,
) {
	// Source -> Write loop
	for {
		var data interface{}

		select {
		case data = <-source:
			break
		case <-stop:
			return
		}

		// Write out
		if write(data) != nil {
			return
		}
	}
}
