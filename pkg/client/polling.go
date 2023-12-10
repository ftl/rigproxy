package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ftl/rigproxy/pkg/protocol"
)

// ResponseHandler is a callback to handle a response.
type ResponseHandler interface {
	Handle(protocol.Response)
}

// ResponseHandlerFunc wraps a function matching the Handle signature to implement the ResponseHandler interface.
type ResponseHandlerFunc func(protocol.Response)

// Handle the given response
func (f ResponseHandlerFunc) Handle(r protocol.Response) {
	f(r)
}

// PollRequest contains a command with arguments that should be send perodically to a rigctld server.
// The given handler is used to handle the responses from the rigctld server.
type PollRequest struct {
	Command protocol.Command
	Args    []string
	Handler ResponseHandler
}

// PollCommand creates a PollRequest from the given handler, command name and arguments.
func PollCommand(handler ResponseHandler, command string, args ...string) PollRequest {
	var cmd protocol.Command
	if len(command) == 1 {
		cmd = protocol.ShortCommand(command)
	} else {
		cmd = protocol.LongCommand(command)
	}
	return PollRequest{
		Command: cmd,
		Args:    args,
		Handler: handler,
	}
}

// PollCommandFunc creates a PollRequest from the given handler function, command name and arguments.
func PollCommandFunc(f func(protocol.Response), command string, args ...string) PollRequest {
	return PollCommand(ResponseHandlerFunc(f), command, args...)
}

type polling struct {
	tick         *time.Ticker
	requestsLock *sync.RWMutex
	requests     []PollRequest
	done         chan struct{}
}

func startPolling(trx *protocol.Transceiver, interval time.Duration, timeout time.Duration, requests []PollRequest) *polling {
	result := polling{
		tick:         time.NewTicker(interval),
		requestsLock: new(sync.RWMutex),
		requests:     requests,
		done:         make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-result.tick.C:
				result.requestsLock.RLock()
				requests := result.requests
				result.requestsLock.RUnlock()

				result.poll(trx, timeout, requests)
			case <-result.done:
				return
			}
		}
	}()

	return &result
}

func (p *polling) stop() {
	select {
	case <-p.done:
	default:
		p.tick.Stop()
		close(p.done)
	}
}

func (p *polling) poll(trx *protocol.Transceiver, timeout time.Duration, requests []PollRequest) {
	for _, pollRequest := range requests {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		request := protocol.Request{Command: pollRequest.Command, Args: pollRequest.Args}
		response, err := trx.Send(ctx, request)
		if err != nil {
			log.Printf("sending poll request %s failed: %v", pollRequest.Command.Long, err)
			if errors.Is(err, protocol.ErrFeatureNotAvailable) || errors.Is(err, protocol.ErrFeatureNotImplemented) || errors.Is(err, protocol.ErrFunctionDeprecated) {
				log.Printf("deactivating poll request %s: %v", pollRequest.Command.Long, err)
				p.remove(pollRequest.Command)
			}
			continue
		}
		if response.Result != "0" {
			log.Printf("poll request %s failed with result: %s", pollRequest.Command.Long, response.Result)
			continue
		}
		pollRequest.Handler.Handle(response)
	}
}

func (p *polling) add(request PollRequest) {
	p.requestsLock.Lock()
	defer p.requestsLock.Unlock()

	for i, pollRequest := range p.requests {
		if pollRequest.Command == request.Command {
			p.requests[i] = request
			return
		}
	}

	p.requests = append(p.requests, request)
}

func (p *polling) remove(command protocol.Command) {
	p.requestsLock.Lock()
	defer p.requestsLock.Unlock()

	for i, pollRequest := range p.requests {
		if pollRequest.Command != command {
			continue
		}

		// order is not relevant, just swap the last element to the free position and cut 1 from the end
		p.requests[i] = p.requests[len(p.requests)-1]
		p.requests = p.requests[:len(p.requests)-1]
		return
	}
}

// StartPolling the connected rigctld server with the given interval and timeout and the given set of requests.
// Poll requests can be added and removed on demand using AddPolls and RemovePolls.
func (c *Conn) StartPolling(interval time.Duration, timeout time.Duration, requests ...PollRequest) error {
	if c.polling != nil {
		return fmt.Errorf("polling is already active")
	}

	c.polling = startPolling(c.trx, interval, timeout, requests)
	return nil
}

// StopPolling stops the polling loop.
func (c *Conn) StopPolling() {
	if c.polling == nil {
		return
	}

	c.polling.stop()
	c.polling = nil
}

// IsPolling indicates if this connection is polling the rigctld server periodically.
func (c *Conn) IsPolling() bool {
	return c.polling != nil
}

// AddPolls while polling is already active. If there is already a poll request with the given command
// in the list of poll requests, the new request replaces the old one.
func (c *Conn) AddPolls(requests ...PollRequest) {
	if c.polling == nil {
		panic("not polling")
	}

	for _, pollRequest := range requests {
		c.polling.add(pollRequest)
	}
}

// Remove the poll requests with the given command from the list of poll requests.
func (c *Conn) RemovePolls(commands ...protocol.Command) {
	if c.polling == nil {
		return
	}

	for _, command := range commands {
		c.polling.remove(command)
	}
}
