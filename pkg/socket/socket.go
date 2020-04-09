package socket

import (
	"log"
	"sync"

	"github.com/bettercallshao/kut/pkg/exch"
	"github.com/bettercallshao/kut/pkg/msg"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ChannelSignal contains signals required for socket exchange
type ChannelSignal struct {
	BackLock  chan int
	FrontLock chan int
	Command   chan interface{}
	Output    chan interface{}
}

// ChannelSignalMap is a map from channel name to signals
type ChannelSignalMap map[string]ChannelSignal

var sigMap ChannelSignalMap
var mux sync.Mutex

// Validate validates if context is ready for socket
func Validate(c *gin.Context) {
	mux.Lock()
	defer mux.Unlock()

	if sigMap == nil {
		sigMap = make(ChannelSignalMap)
	}

	channel := c.Param("channel")
	if channel != "" {
		_, ok := sigMap[channel]
		if !ok {
			sigMap[channel] = ChannelSignal{
				make(chan int, 1),
				make(chan int, 1),
				make(chan interface{}),
				// Allow output channel to buffer msgs
				make(chan interface{}, 1024),
			}
		}
	}
}

// FrontGET handles get requests from browser
func FrontGET(c *gin.Context) {
	read, write, close, err := createFunctors(c, false)
	if err != nil {
		return
	}
	defer close()

	log.Println("starting exchange loop - frontend ...")
	channel := c.Param("channel")
	exch.Exchange(
		sigMap[channel].Output,
		sigMap[channel].Command,
		read,
		write,
		sigMap[channel].FrontLock,
		nil,
	)
	log.Println("finished exchange loop - frontend")
}

// BackGET handles get requests from executor
func BackGET(c *gin.Context, coInit func(bool) error) {
	read, write, close, err := createFunctors(c, true)
	if err != nil {
		return
	}
	defer close()

	log.Println("starting exchange loop - backend ...")
	channel := c.Param("channel")
	exch.Exchange(
		sigMap[channel].Command,
		sigMap[channel].Output,
		read,
		write,
		sigMap[channel].BackLock,
		coInit,
	)
	log.Println("finished exchange loop - backend")
}

func createFunctors(
	c *gin.Context,
	isOutput bool,
) (
	func() (interface{}, error),
	func(interface{}) error,
	func(),
	error,
) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(msg.HBAD, msg.Error{Error: "Failed to set up websocket."})
		log.Println("ws upgrade failed")
		return nil, nil, nil, err
	}
	log.Println("opened ws connection")

	read := func() (interface{}, error) {
		data := msg.Command{}
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println("ws read failed")
		}
		return data, err
	}

	if isOutput {
		read = func() (interface{}, error) {
			data := msg.Output{}
			err := conn.ReadJSON(&data)
			if err != nil {
				log.Println("ws read failed")
			}
			return data, err
		}
	}

	write := func(data interface{}) error {
		err := conn.WriteJSON(data)
		if err != nil {
			log.Println("ws write failed")
		}
		return err
	}

	close := func() {
		log.Println("closing ws connection ...")
		conn.Close()
	}

	return read, write, close, nil
}
