package protocol

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

type Transceiver struct {
	rw       io.ReadWriter
	outgoing chan transmission
	polling  polling
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
		outgoing: make(chan transmission, 1),
		polling: polling{
			tick: time.NewTicker(1 * time.Second),
		},
		closed: make(chan struct{}),
	}
	result.polling.tick.Stop()

	go result.start()

	return &result
}

func NewPollingTransceiver(rw io.ReadWriter, interval time.Duration, timeout time.Duration, requests ...PollRequest) *Transceiver {
	result := Transceiver{
		rw:       rw,
		outgoing: make(chan transmission),
		polling: polling{
			tick:     time.NewTicker(interval),
			timeout:  timeout,
			requests: requests,
		},
		closed: make(chan struct{}),
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
			} else if resp.Result != "0" {
				log.Printf("hamlib error code: %s", resp.Result)
				tx.err <- fmt.Errorf("request failed: %s", resp.Result)
			} else {
				select {
				case tx.response <- resp:
				default:
					log.Printf("could not queue response to transmission, nobody is listening: %+v", resp)
				}
			}
		case <-t.polling.tick.C:
			go t.poll()
		}
	}
}

func (t *Transceiver) Send(ctx context.Context, req Request) (Response, error) {
	select {
	case <-t.closed:
		return Response{}, errors.New("transceiver already closed")
	default:
	}

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

func (t *Transceiver) Close() {
	select {
	case <-t.closed:
		return
	default:
		t.polling.tick.Stop()
		close(t.closed)
	}
}

func (t *Transceiver) WhenDone(f func()) {
	go func() {
		<-t.closed
		f()
	}()
}

type polling struct {
	tick     *time.Ticker
	timeout  time.Duration
	requests []PollRequest
}

type ResponseHandler interface {
	Handle(Request, Response) error
}

type ResponseHandlerFunc func(Request, Response) error

func (f ResponseHandlerFunc) Handle(request Request, response Response) error {
	return f(request, response)
}

type PollRequest struct {
	Command Command
	Args    []string
	Handler ResponseHandler
}

func PollCommandFunc(f func(Request, Response) error, command string, args ...string) PollRequest {
	return PollCommand(ResponseHandlerFunc(f), command, args...)
}

func PollCommand(handler ResponseHandler, command string, args ...string) PollRequest {
	var cmd Command
	if len(command) == 1 {
		cmd = ShortCommand(command)
	} else {
		cmd = LongCommand(command)
	}
	return PollRequest{
		Command: cmd,
		Args:    args,
		Handler: handler,
	}
}

func (t Transceiver) poll() {
	for _, r := range t.polling.requests {
		ctx, _ := context.WithTimeout(context.Background(), t.polling.timeout)
		request := Request{Command: r.Command, Args: r.Args}
		response, err := t.Send(ctx, request)
		if err != nil {
			log.Printf("sending poll request %s failed: %v", r.Command.Long, err)
			continue
		}

		err = r.Handler.Handle(request, response)
		if err != nil {
			log.Printf("receiving poll response %s failed: %v", r.Command.Long, err)
		}
	}
}
