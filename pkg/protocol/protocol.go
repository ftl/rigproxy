package protocol

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type CommandKey string

const NoCommand = CommandKey("")

func subCommandKey(cmd string, sub string) CommandKey {
	return CommandKey(cmd + "_" + sub)
}

type Command struct {
	Short                string
	Long                 string
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

type Response struct {
	Data   []string
	Result string
}

func (r *Response) Format() string {
	if len(r.Data) > 0 {
		return fmt.Sprintf("%s\nRPRT %s", strings.Join(r.Data, "\n"), r.Result)
	}
	return fmt.Sprintf("RPRT %s", r.Result)
}

func (r *Response) ExtendedFormat(separator string, keys []string) string {
	buffer := bytes.NewBufferString("")

	for i, value := range r.Data {
		fmt.Fprintf(buffer, "%s: %s\n", keys[i], value)
	}
	fmt.Fprintf(buffer, "RPRT %s", r.Result)

	return buffer.String()
}

type Transceiver struct {
	rw       io.ReadWriter
	outgoing chan transmission
	closed   chan struct{}
}

type transmission struct {
	request  Request
	response chan Response
}

func NewTransceiver(rw io.ReadWriter) *Transceiver {
	result := Transceiver{
		rw:       rw,
		outgoing: make(chan transmission, 20),
		closed:   make(chan struct{}),
	}

	go result.start()

	return &result
}

func (t *Transceiver) start() {
	txError := Response{Result: "501"}
	rxError := Response{Result: "502"}
	r := NewResponseReader(t.rw)
	for {
		select {
		case <-t.closed:
			return
		case tx := <-t.outgoing:
			_, err := fmt.Fprintln(t.rw, tx.request.LongFormat())
			if err != nil {
				log.Println("transmit:", err)
				tx.response <- txError
			}
			resp, err := r.ReadResponse()
			if err != nil {
				log.Println("receive:", err)
				tx.response <- rxError
			} else {
				tx.response <- resp
			}
		}
	}
}

func (t *Transceiver) Send(ctx context.Context, req Request) (Response, error) {
	select {
	case <-t.closed:
		return Response{}, errors.New("transceiver already closed")
	default:
		tx := transmission{request: req, response: make(chan Response)}
		t.outgoing <- tx
		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		case resp := <-tx.response:
			return resp, nil
		}
	}
}

func (t *Transceiver) Close() {
	select {
	case <-t.closed:
		return
	default:
		close(t.closed)
	}
}
