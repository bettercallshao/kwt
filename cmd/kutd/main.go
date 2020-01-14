package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bettercallshao/kut/pkg/channel"
	"github.com/bettercallshao/kut/pkg/socket"
)

func main() {
	log.SetPrefix("[kutd] ")
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(channel.Validate)
	router.Use(socket.Validate)

	router.GET("/assets/*filename", assetGET)
	router.NoRoute(indexGET)

	router.GET("/menu", menuGET)
	router.POST("/menu", menuPOST)

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
	router.Run(":" + port)
}
