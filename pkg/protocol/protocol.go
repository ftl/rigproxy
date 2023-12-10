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

type Error struct {
	code    string
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("hamlib error %s: %s", e.code, e.message)
}

func newError(code string) Error {
	message, ok := errorMessagesByCode[code]
	if !ok {
		message = "Unknown error"
	}
	return Error{code: code, message: message}
}

var errorMessagesByCode = map[string]string{
	"-1":  "Invalid parameter",
	"-2":  "Invalid configuration",
	"-3":  "Memory shortage",
	"-4":  "Feature not implemented",
	"-5":  "Communication timed out",
	"-6":  "IO error",
	"-7":  "Internal Hamlib error",
	"-8":  "Protocol error",
	"-9":  "Command rejected by the rig",
	"-10": "Command performed, but arg truncated, result not guaranteed",
	"-11": "Feature not available",
	"-12": "Target VFO unaccessible",
	"-13": "Communication bus error",
	"-14": "Communication bus collision",
	"-15": "NULL RIG handle or invalid pointer parameter",
	"-16": "Invalid VFO",
	"-17": "Argument out of domain of func",
	"-18": "Function deprecated",
	"-19": "Security error password not provided or crypto failure",
	"-20": "Rig is not powered on",
}

var (
	ErrInvalidParameter          = newError("-1")
	ErrInvalidConfiguration      = newError("-2")
	ErrMemoryShortage            = newError("-3")
	ErrFeatureNotImplemented     = newError("-4")
	ErrCommunicationTimedOut     = newError("-5")
	ErrIOError                   = newError("-6")
	ErrInternalHamlibError       = newError("-7")
	ErrProtocolError             = newError("-8")
	ErrCommandRejectedByRig      = newError("-9")
	ErrArgTruncated              = newError("-10")
	ErrFeatureNotAvailable       = newError("-11")
	ErrTargetVFOUnaccessible     = newError("-12")
	ErrCommunicationBusError     = newError("-13")
	ErrCommunicationBusCollision = newError("-14")
	ErrNullRigHandle             = newError("-15")
	ErrInvalidVFO                = newError("-16")
	ErrArgumentOutOfDomain       = newError("-17")
	ErrFunctionDeprecated        = newError("-18")
	ErrSecurityError             = newError("-19")
	ErrRigNotPoweredOn           = newError("-20")
)
