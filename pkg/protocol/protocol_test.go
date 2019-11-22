package protocol

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ftl/rigproxy/pkg/test"
)

func TestSimpleCommandKey(t *testing.T) {
	cmd := Command{Short: 'a', Long: "get_a"}
	req := Request{Command: cmd}
	assert.Equal(t, CommandKey("get_a"), req.Key())
}

func TestCommandKeyWithSubCommand(t *testing.T) {
	cmd := Command{Short: 'b', Long: "get_b", HasSubCommand: true}
	req := Request{Command: cmd, Args: []string{"first"}}
	assert.Equal(t, CommandKey("get_b_first"), req.Key())
}

func TestInvalidatingCommandKey(t *testing.T) {
	cmd := Command{Short: 'c', Long: "set_c", InvalidatesCommand: "get_c"}
	req := Request{Command: cmd}
	assert.Equal(t, CommandKey("get_c"), req.InvalidatedKey())
}

func TestInvalidatingCommandKeyWithSubCommand(t *testing.T) {
	cmd := Command{Short: 'd', Long: "set_d", InvalidatesCommand: "get_d", HasSubCommand: true}
	req := Request{Command: cmd, Args: []string{"first"}}
	assert.Equal(t, CommandKey("get_d_first"), req.InvalidatedKey())
}

func TestTransceiverSendReceiveRoundtrip(t *testing.T) {
	buffer := test.NewBuffer("get_freq:\nFrequency: 3720000\nRPRT 0\nRPRT 11\n")

	trx := NewTransceiver(buffer)
	defer trx.Close()

	resp, err := trx.Send(context.Background(), Request{Command: ShortCommand("f")})
	assert.NoError(t, err)
	assert.Equal(t, Response{
		Command: CommandKey("get_freq"),
		Data:    []string{"3720000"},
		Keys:    []string{"Frequency"},
		Result:  "0",
	}, resp)

	resp, err = trx.Send(context.Background(), Request{Command: ShortCommand("f")})
	assert.NoError(t, err)
	assert.Equal(t, Response{
		Result: "11",
	}, resp)

	buffer.AssertWritten(t, "+\\get_freq\n+\\get_freq\n")
}
