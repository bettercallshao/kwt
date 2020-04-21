package msg

import (
	"github.com/bettercallshao/kut/pkg/cmd"
	"github.com/bettercallshao/kut/pkg/menu"
)

// HGOOD is a good code for http calls
const HGOOD = 200

// HBAD is a bad code for http calls
const HBAD = 400

// MenuIngest contains info for ingesting menu.
type MenuIngest struct {
	Source string `json:"source" binding:"required"`
}

// ChannelInfo contains info for operating a channel
type ChannelInfo struct {
	Name  string    `json:"name" binding:"required"`
	Token string    `json:"token" binding:"-"`
	Menu  menu.Menu `json:"menu" binding:"-"`
}

// ChannelInfoMap maps channel name to infos
type ChannelInfoMap map[string]ChannelInfo

// MenuName is an alias for the json field name
const MenuName = "menuname"

// Version represents kut version
type Version struct {
	Version string `json:"version"`
}

// Error for http requests
type Error struct {
	Error string `json:"error"`
}

// List for http response with list
type List struct {
	List interface{} `json:"list"`
}

// Empty for http success with no data
type Empty struct{}

// Command for user submission
type Command struct {
	Token  string      `json:"token" binding:"required"`
	Action menu.Action `json:"action" binding:"required"`
}

// Output from command
type Output struct {
	Token   string      `json:"token" binding:"required"`
	Payload cmd.Payload `json:"payload" binding:"required"`
}
