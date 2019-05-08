package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	for _, cmd := range Commands {
		assert.Equal(t, cmd, ShortCommands[cmd.Short])
		assert.Equal(t, cmd, LongCommands[cmd.Long])
	}
}
