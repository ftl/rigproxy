package protocol

import (
	"bytes"
	"fmt"
	"strings"
)

type CommandKey string

const NoCommand = CommandKey("")

func subCommandKey(cmd string, sub string) CommandKey {
	return CommandKey(cmd + "_" + sub)
}

type Command struct {
	Short                byte
	Long                 string
	Args                 int
	ArgsInLine           bool
	InvalidatesCommand   string
	HasSubCommand        bool
	SupportsExtendedMode bool
	Cacheable            bool
}

type Request struct {
	Command
	ExtendedSeparator string
	Args              []string
}

func (r *Request) Key() CommandKey {
	if r.HasSubCommand && len(r.Args) > 0 {
		return subCommandKey(r.Long, r.Args[0])
	}
	return CommandKey(r.Long)
}

func (r *Request) InvalidatedKey() CommandKey {
	if r.InvalidatesCommand != "" {
		if r.HasSubCommand && len(r.Args) > 0 {
			return subCommandKey(r.InvalidatesCommand, r.Args[0])
		}
		return CommandKey(r.InvalidatesCommand)
	}
	return NoCommand
}

func (r *Request) LongFormat() string {
	return strings.Join(append([]string{"\\" + r.Long}, r.Args...), " ")
}

func (r *Request) ExtendedFormat() string {
	return "+" + r.LongFormat()
}

type Response struct {
	Command CommandKey
	Data    []string
	Keys    []string
	Result  string
}

func (r *Response) Format() string {
	if len(r.Data) == 0 || r.Result != "0" {
		return fmt.Sprintf("RPRT %s", r.Result)
	}
	return strings.Join(r.Data, "\n")
}

func (r *Response) ExtendedFormat(separator string) string {
	buffer := bytes.NewBufferString("")

	fmt.Fprintf(buffer, "%s:%s", r.Command, separator)
	for i, value := range r.Data {
		if r.Keys[i] != "" {
			fmt.Fprintf(buffer, "%s: %s%s", r.Keys[i], value, separator)
		} else {
			fmt.Fprintf(buffer, "%s%s", value, separator)
		}
	}
	fmt.Fprintf(buffer, "RPRT %s", r.Result)

	return buffer.String()
}
