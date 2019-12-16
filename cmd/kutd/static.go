package main

import (
	"io/ioutil"
	"strings"

	"github.com/bettercallshao/kut/pkg/msg"
	"github.com/gin-gonic/gin"
)

func assetWrite(c *gin.Context, filename string) {
	file, ok := Assets.Files[strings.TrimLeft(filename, "/")]
	if !ok {
		c.String(msg.HBAD, "")
		return
	}

	file.Seek(0, 0)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(msg.HBAD, "")
		return
	}

	cType := "text/html"
	if strings.HasSuffix(filename, ".css") {
		cType = "text/css"
	}
	c.Data(msg.HGOOD, cType, data)
}

func indexGET(c *gin.Context) {
	assetWrite(c, "index.html")
}

func assetGET(c *gin.Context) {
	assetWrite(c, c.Param("filename"))
}
