package channel

import (
	"errors"
	"log"
	"sync"

	"github.com/bettercallshao/kut/pkg/menu"
	"github.com/bettercallshao/kut/pkg/msg"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

var mux sync.Mutex
var infoMap msg.ChannelInfoMap

// List returns a list of channels
func List() []string {
	return []string{"0", "1", "2"}
}

// Validate validates if url's channel is good
func Validate(c *gin.Context) {
	// Sneak in a hack here to initialize map
	mux.Lock()
	if infoMap == nil {
		infoMap = make(msg.ChannelInfoMap)
	}
	mux.Unlock()

	channel := c.Param("channel")
	if channel != "" {
		found := false
		for _, good := range List() {
			if channel == good {
				found = true
				break
			}
		}
		if !found {
			c.JSON(msg.HBAD, msg.Error{Error: "Invalid channel."})
			c.Abort()
		}
	}
}

// ListGET responds with a list of channels
func ListGET(c *gin.Context) {
	c.JSON(msg.HGOOD, msg.List{List: List()})
}

// GET handles get requests
func GET(c *gin.Context) {
	mux.Lock()
	defer mux.Unlock()

	var channel string
	var ok bool
	var info msg.ChannelInfo

	channel = c.Param("channel")
	info, ok = infoMap[channel]
	if !ok {
		c.JSON(msg.HBAD, msg.Error{Error: "Channel is idle."})
		return
	}

	c.JSON(msg.HGOOD, info)
}

// CreateInfo creates the info object from context
func CreateInfo(c *gin.Context) (msg.ChannelInfo, error) {
	info := msg.ChannelInfo{}
	err := errors.New("")

	menuName := c.Query(msg.MenuName)
	if menuName == "" {
		c.JSON(msg.HBAD, msg.Error{Error: "Invalid query param."})
		log.Println("invalid query param")
		return info, err
	}

	info.Menu, err = menu.Load(menuName)
	if err != nil {
		c.JSON(msg.HBAD, msg.Error{Error: "Failed to load menu."})
		log.Println("failed to load menu")
		return info, err
	}

	info.Name = c.Param("channel")
	info.Token = randstr.String(4)

	return info, nil
}

// AddInfo adds info to map
func AddInfo(info msg.ChannelInfo) {
	mux.Lock()
	defer mux.Unlock()

	log.Printf("adding info key: %s <= token: %s", info.Name, info.Token)
	infoMap[info.Name] = info
}

// DeleteInfo deletes info from map
func DeleteInfo(info msg.ChannelInfo) {
	mux.Lock()
	defer mux.Unlock()

	log.Printf("deleting info key: %s", info.Name)
	delete(infoMap, info.Name)
}
