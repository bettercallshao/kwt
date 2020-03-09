package exch

import (
	"errors"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestHappyPath(t *testing.T) {
	sema := make(smChan, 1)
	sourceSink := make(ifChan)

	readCount := 10
	read := func() (interface{}, error) {
		time.Sleep(5 * time.Millisecond)
		if readCount > 0 {
			readCount--
		} else {
			return readCount, errors.New("")
		}

		return readCount, nil
	}

	writeCount := 10
	write := func(interface{}) error {
		writeCount--
		return nil
	}

	coInitCount := 2
	coInit := func(bool) error {
		coInitCount--
		return nil
	}

	err := Exchange(sourceSink, sourceSink, read, write, sema, coInit)
	assert.NilError(t, err)
	assert.Equal(t, readCount, 0)
	assert.Equal(t, writeCount, 0)
	assert.Equal(t, coInitCount, 0)
}
