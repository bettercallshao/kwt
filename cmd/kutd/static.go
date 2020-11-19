package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/bettercallshao/kut/pkg/msg"
	"github.com/gin-gonic/gin"
)

func readAsset(filename string) ([]byte, bool) {
	file, ok := Assets.Files[filename]
	if !ok {
		return nil, false
	}

	file.Seek(0, 0)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false
	}

	return data, true
}

func readFile(filename string, assetRoot string) ([]byte, bool) {
	data, err := ioutil.ReadFile(path.Join(assetRoot, filename))
	return data, err == nil
}

func staticWrite(c *gin.Context, filename string) {
	assetRoot := os.Getenv("ASSETS_ROOT")
	var data []byte
	var ok bool
	if assetRoot == "" {
		data, ok = readAsset(filename)
	} else {
		// TODO: vulnerable to relative path attack
		data, ok = readFile(filename, assetRoot)
	}

	if !ok {
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
	staticWrite(c, "index.html")
}

func assetGET(c *gin.Context) {
	staticWrite(c, strings.TrimLeft(c.Param("filename"), "/"))
}
