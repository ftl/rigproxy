package protocol

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

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
	conn     net.Conn
	outgoing chan transmission
	closed   chan struct{}
	timeout  time.Duration
}

type transmission struct {
	request  Request
	response chan Response
}

func NewTransceiver(conn net.Conn) *Transceiver {
	result := Transceiver{
		conn:     conn,
		outgoing: make(chan transmission),
		closed:   make(chan struct{}),
		timeout:  500 * time.Millisecond,
	}

	go result.start()

	return &result
}

func (t *Transceiver) start() {
	txError := Response{Result: "501"}
	rxError := Response{Result: "502"}
	connectionClosed := Response{Result: "503"}
	r := NewResponseReader(t.conn, t.timeout)
	for {
		select {
		case <-t.closed:
			return
		case tx := <-t.outgoing:
			_, err := fmt.Fprintln(t.conn, tx.request.ExtendedFormat())
			if err != nil {
				log.Println("transmit:", err)
				tx.response <- txError
			}
			resp, err := r.ReadResponse(tx.request.SupportsExtendedMode)
			if err == io.EOF {
				log.Println("receive: connection closed")
				tx.response <- connectionClosed
				close(t.closed)
				return
			} else if err != nil {
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
		resp := <-tx.response
		return resp, nil
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
