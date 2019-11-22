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
	Short                byte
	Long                 string
	Args                 int
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

	fmt.Fprintf(buffer, "%s:\n", r.Command)
	for i, value := range r.Data {
		if r.Keys[i] != "" {
			fmt.Fprintf(buffer, "%s: %s\n", r.Keys[i], value)
		} else {
			fmt.Fprintln(buffer, value)
		}
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
	err      chan error
}

func NewTransceiver(rw io.ReadWriter) *Transceiver {
	result := Transceiver{
		rw:       rw,
		outgoing: make(chan transmission),
		closed:   make(chan struct{}),
	}

	go result.start()

	return &result
}

func (t *Transceiver) start() {
	r := NewResponseReader(t.rw)
	for {
		select {
		case <-t.closed:
			return
		case tx := <-t.outgoing:
			_, err := fmt.Fprintln(t.rw, tx.request.ExtendedFormat())
			if err != nil {
				log.Println("transmit:", err)
				tx.err <- fmt.Errorf("transmission of request failed: %w", err)
			}
			resp, err := r.ReadResponse(tx.request.SupportsExtendedMode)
			if err == io.EOF {
				log.Println("receive: connection closed")
				tx.err <- fmt.Errorf("connection closed while waiting for response: %w", err)
				close(t.closed)
				return
			} else if err != nil {
				log.Println("receive:", err)
				tx.err <- fmt.Errorf("receiving of response failed: %w", err)
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
		tx := transmission{request: req, response: make(chan Response), err: make(chan error)}
		t.outgoing <- tx
		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		case err := <-tx.err:
			return Response{}, err
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

func (t *Transceiver) WhenDone(f func()) {
	go func() {
		<-t.closed
		f()
	}()
}
