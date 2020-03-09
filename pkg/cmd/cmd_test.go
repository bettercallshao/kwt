package cmd

import (
	"testing"

	"gotest.tools/assert"
)

func TestHappyPath(t *testing.T) {
	sink := make(chan Payload)
	go Run("echo abc", sink)

	payload := <-sink
	assert.Equal(t, payload.Text, ">> echo abc")

	payload = <-sink
	assert.Equal(t, payload.Text, "abc")
}
