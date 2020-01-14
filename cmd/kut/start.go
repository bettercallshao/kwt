package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/bettercallshao/kut/pkg/cmd"
	"github.com/bettercallshao/kut/pkg/exch"
	"github.com/bettercallshao/kut/pkg/menu"
	"github.com/bettercallshao/kut/pkg/msg"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

func getInfo(master string, channel string) (msg.ChannelInfo, error) {
	info := msg.ChannelInfo{}
	herr := msg.Error{}
	resp, err := resty.New().R().
		SetResult(&info).
		SetError(&herr).
		Get(join(master, "channel", channel))
	if err != nil {
		return info, err
	}
	if !resp.IsSuccess() {
		return info, errors.New("")
	}
	return info, nil
}

func start(master string, channel string, menuName string) error {
	var conn *websocket.Conn
	var info msg.ChannelInfo
	var err error

	// Construct url
	url := addParam(
		join(wsURL(master), "ws", "back", channel),
		msg.MenuName,
		menuName)
	log.Printf("connecting to master %s ...\n", url)

	// Open ws connection
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get channel info
	time.Sleep(time.Millisecond)
	info, err = getInfo(master, channel)
	if err != nil {
		return err
	}
	log.Printf("channel definition:\n%s\n", jsonParagraph(info))

	// Set up a cancel channel
	cancelMain := cancelChan()
	cancelProc := cancelChan()

	// Start exchange
	source := make(chan interface{})
	sink := make(chan interface{})
	read := func() (interface{}, error) {
		data := msg.Command{}
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println("ws read failed")
		}
		return data, err
	}
	write := func(data interface{}) error {
		err := conn.WriteJSON(data)
		if err != nil {
			log.Println("ws write failed")
		}
		return err
	}
	coInit := func(start bool) error {
		if start {
			log.Println("starting exchange loop ...")
		} else {
			log.Println("finished exchange loop")
			log.Println("signaling processor termination ...")
			cancelProc <- os.Interrupt
			cancelMain <- os.Interrupt
		}
		return nil
	}
	go exch.Exchange(
		source,
		sink,
		read,
		write,
		nil,
		coInit,
	)

	// Start processor loop
	go func() {
		log.Println("starting processor loop ...")
		for {
			var raw interface{}

			select {
			case raw = <-sink:
				break
			case <-cancelProc:
				log.Println("finished processor loop")
				return
			}

			command, ok := raw.(msg.Command)
			if !ok {
				log.Printf("invalid command: %s\n", raw)
				continue
			}

			log.Printf("command received:\n%s\n", jsonParagraph(command))

			if !ok || command.Token != info.Token {
				log.Printf("token mismatch: %s\n", command.Token)
				continue
			}

			input, err := menu.Render(command.Action)
			if err != nil {
				log.Println("action failed to render")
				continue
			}

			cmdSink := make(chan cmd.Payload)
			log.Printf("running command: %s\n", input)
			go cmd.Run(input, cmdSink)

			for payload := range cmdSink {
				log.Printf("sending payload:\n%s\n", jsonParagraph(payload))
				source <- msg.Output{
					Token:   info.Token,
					Payload: payload,
				}
			}
		}
	}()

	// Trap until user cancel
	<-cancelMain
	return nil
}
