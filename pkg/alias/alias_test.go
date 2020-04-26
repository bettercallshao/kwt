package alias

import (
	"testing"

	"gotest.tools/assert"
)

func TestHappyPath(t *testing.T) {
	s := New()
	Avoid(s, []string{"a"})

	assert.Equal(t, Pick(s, "start")[0], "s")
	assert.Equal(t, Pick(s, "stop")[0], "t")
	assert.Equal(t, len(Pick(s, "sta")), 0)
	assert.Equal(t, Pick(s, " -,.x")[0], "x")
}
