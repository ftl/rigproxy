package protocol

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {
	testCases := []struct {
		desc     string
		value    string
		valid    bool
		expected Request
	}{
		{"read frequency short", "f", true, Request{Command: ShortCommand("f"), Args: []string{}}},
		{"read frequency long", "\\get_freq", true, Request{Command: LongCommand("get_freq"), Args: []string{}}},
		{"read frequency extended short", "+f", true, Request{Command: ShortCommand("f"), ExtendedSeparator: "\n", Args: []string{}}},
		{"write frequency short", "F 14074000", true, Request{Command: ShortCommand("F"), Args: []string{"14074000"}}},
		{"write frequency long", "\\set_freq 3720000", true, Request{Command: LongCommand("set_freq"), Args: []string{"3720000"}}},
		{"write frequency long extended", ";\\set_freq 3720000", true, Request{Command: LongCommand("set_freq"), ExtendedSeparator: ";", Args: []string{"3720000"}}},
		{"get functions short", "u ?", true, Request{Command: ShortCommand("u"), Args: []string{"?"}}},
		{"get functions long", "\\get_func ?", true, Request{Command: LongCommand("get_func"), Args: []string{"?"}}},
		{"get functions extended long", ",\\get_func ?", true, Request{Command: LongCommand("get_func"), Args: []string{"?"}}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := ParseRequest(tC.value)
			if tC.valid {
				assert.NoError(t, err)
				assert.Equal(t, tC.expected, actual)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestRequestReader(t *testing.T) {
	buffer := bytes.NewBufferString("# A comment before the command\nF 140740000\n# A comment after the command\n")

	reader := NewRequestReader(buffer)
	req, err := reader.ReadRequest()

	require.NoError(t, err)
	assert.Equal(t, Request{
		Command: ShortCommand("F"),
		Args:    []string{"140740000"},
	}, req)

	_, err = reader.ReadRequest()
	assert.Equal(t, io.EOF, err)
}

func TestResponseReader(t *testing.T) {
	buffer := bytes.NewBufferString("USB\n2400\nRPRT 0\nget_freq:\nFrequency: 145000000\nRPRT 0\nRPRT 11\n")

	reader := NewResponseReader(buffer)

	resp, err := reader.ReadResponse(false)
	require.NoError(t, err)
	assert.Equal(t, Response{
		Data:   []string{"USB", "2400"},
		Result: "0",
	}, resp)

	resp, err = reader.ReadResponse(true)
	require.NoError(t, err)
	assert.Equal(t, Response{
		Command: CommandKey("get_freq"),
		Data:    []string{"145000000"},
		Keys:    []string{"Frequency"},
		Result:  "0",
	}, resp)

	resp, err = reader.ReadResponse(false)
	assert.NoError(t, err)
	assert.Equal(t, resp.Result, "11")

	_, err = reader.ReadResponse(false)
	assert.Equal(t, io.EOF, err)
}
