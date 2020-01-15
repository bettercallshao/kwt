package main

import (
	"github.com/bettercallshao/kut/pkg/menu"
	"github.com/bettercallshao/kut/pkg/msg"
	"github.com/gin-gonic/gin"
)

func menuGET(c *gin.Context) {
	c.JSON(msg.HGOOD, msg.List{List: menu.List()})
}

func menuPOST(c *gin.Context) {
	var err error
	var ingest msg.MenuIngest

	if c.BindJSON(&ingest) != nil {
		c.JSON(msg.HBAD, msg.Error{Error: "Invalid JSON input."})
		return
	}

	if ingest.Name == "" || ingest.Source == "" {
		c.JSON(msg.HBAD, msg.Error{Error: "Name and source cannot be empty."})
		return
	}

	if err = menu.Ingest(ingest.Name, ingest.Source); err != nil {
		c.JSON(msg.HBAD, msg.Error{Error: "Failed to ingest."})
		return
	}

	c.JSON(msg.HGOOD, msg.Empty{})
}

func menuItemGET(c *gin.Context) {
	menu, err := menu.Load(c.Param("menu"))
	if err != nil {
		c.JSON(msg.HBAD, msg.Error{Error: "Failed to load menu."})
		return
	}
	c.JSON(msg.HGOOD, menu)
}
