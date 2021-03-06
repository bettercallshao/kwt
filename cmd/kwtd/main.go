package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bettercallshao/kwt/pkg/channel"
	"github.com/bettercallshao/kwt/pkg/msg"
	"github.com/bettercallshao/kwt/pkg/socket"
	"github.com/bettercallshao/kwt/pkg/version"
)

func main() {
	log.SetPrefix("[kwtd] ")
	log.Printf("version: %s", version.Version)
	log.Println("starting kwtd ...")

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.Default()
	router.Use(channel.Validate)
	router.Use(socket.Validate)

	router.GET("/assets/*filename", assetGET)
	router.NoRoute(indexGET)

	router.GET("/version", func(c *gin.Context) {
		c.JSON(msg.HGOOD, msg.Version{Version: version.Version})
	})

	router.GET("/menu", menuGET)
	router.POST("/menu", menuPOST)
	router.GET("/menu/:menu", menuItemGET)

	router.GET("/channel", channel.ListGET)
	router.GET("/channel/:channel", channel.GET)

	router.GET("/ws/front/:channel", socket.FrontGET)
	router.GET("/ws/back/:channel", func(c *gin.Context) {
		info, err := channel.CreateInfo(c)
		if err != nil {
			return
		}

		socket.BackGET(c, func(start bool) error {
			if start {
				channel.AddInfo(info)
			} else {
				channel.DeleteInfo(info)
			}
			return nil
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "7171"
	}
	host := "127.0.0.1"
	log.Println("listening on http://" + host + ":" + port)
	router.Run(host + ":" + port)
}
