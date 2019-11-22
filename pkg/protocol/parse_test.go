package protocol

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestReader(t *testing.T) {
	buffer := bytes.NewBufferString(`# A comment before the command
F 14074000
fmv
+\set_mode PKTUSB 1800 # switch to data mode
# A comment after the command
`)
	expectedRequests := []Request{
		{Command: ShortCommand("F"), Args: []string{"14074000"}},
		{Command: ShortCommand("f")},
		{Command: ShortCommand("m")},
		{Command: ShortCommand("v")},
		{Command: LongCommand("set_mode"), Args: []string{"PKTUSB", "1800"}, ExtendedSeparator: "\n"},
	}
	reader := NewRequestReader(buffer)

	for i, expected := range expectedRequests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			req, err := reader.ReadRequest()

			assert.NoError(t, err)
			assert.Equal(t, expected, req)

		})
	}

	_, err := reader.ReadRequest()
	assert.Equal(t, io.EOF, err)
}

func TestEmptyRequestReader(t *testing.T) {
	buffer := bytes.NewBufferString("")
	reader := NewRequestReader(buffer)
	_, err := reader.ReadRequest()
	assert.Equal(t, io.EOF, err)
}

func TestLineBreakIsSpace(t *testing.T) {
	assert.True(t, unicode.IsSpace('\n'))
}

func TestNextRequest(t *testing.T) {
	testCases := []struct {
		desc     string
		value    string
		expected Request
		valid    bool
	}{
		{"empty", "", Request{}, false},
		{"comment", " # a comment", Request{}, false},
		{"single short command", "f", Request{Command: ShortCommand("f")}, true},
		{"multiple short commands", "fmv", Request{Command: ShortCommand("f")}, true},
		{"short command with args", "F 14074000", Request{Command: ShortCommand("F"), Args: []string{"14074000"}}, true},
		{"single long command with args", "\\set_freq 3720000", Request{Command: LongCommand("set_freq"), Args: []string{"3720000"}}, true},
		{"extended long command", ";\\get_freq", Request{Command: LongCommand("get_freq"), ExtendedSeparator: ";"}, true},
		{"extended long command newline", "+\\get_mode", Request{Command: LongCommand("get_mode"), ExtendedSeparator: "\n"}, true},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tC.value)
			actual, err := nextRequest(buffer)
			if tC.valid {
				assert.NoError(t, err)
				assert.Equal(t, tC.expected, actual)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestReadLongCommand(t *testing.T) {
	testCases := []struct {
		desc     string
		value    string
		expected Command
		valid    bool
	}{
		{"empty", "", Command{}, false},
		{"get_freq", "get_freq", LongCommand("get_freq"), true},
		{"  get_freq  ", "  get_freq  ", LongCommand("get_freq"), true},
		{"unknown command", "blah", Command{}, false},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tC.value)
			actual, err := readLongCommand(buffer)
			if tC.valid {
				assert.NoError(t, err)
				assert.Equal(t, tC.expected, actual)
			} else {
				assert.Error(t, err)
			}
		})
	}
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
