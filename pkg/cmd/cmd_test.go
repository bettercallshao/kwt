package cmd

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestHappyPath(t *testing.T) {
	sink := make(chan Payload)
	go Run("echo abc", sink)

	text := ""
	for payload := range sink {
		text += payload.Text
	}
	text = strings.TrimSpace(text)

	assert.Equal(t, text, "abc")
}
