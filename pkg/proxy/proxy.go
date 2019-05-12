package proxy

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ftl/rigproxy/pkg/protocol"
)

type Proxy struct {
	rwc    io.ReadWriteCloser
	trx    Transceiver
	cache  Cache
	closed chan struct{}
}

type Transceiver interface {
	Send(context.Context, protocol.Request) (protocol.Response, error)
}

type Cache interface {
	Put(protocol.CommandKey, protocol.Response)
	Get(protocol.CommandKey) (protocol.Response, bool)
	Invalidate(protocol.CommandKey)
}

var ChkVfoResponse = protocol.Response{
	Command: protocol.CommandKey("chk_vfo"),
	Data:    []string{"CHKVFO 0"},
	Keys:    []string{""},
	Result:  "0",
}

func New(rwc io.ReadWriteCloser, trx Transceiver, done <-chan struct{}) *Proxy {
	return NewCached(rwc, trx, new(nopCache), done)
}

func NewCached(rwc io.ReadWriteCloser, trx Transceiver, cache Cache, done <-chan struct{}) *Proxy {
	result := Proxy{
		rwc:    rwc,
		trx:    trx,
		cache:  cache,
		closed: make(chan struct{}),
	}

	go result.start()
	go func() {
		select {
		case <-done:
			close(result.closed)
			rwc.Close()
		case <-result.closed:
			rwc.Close()
		}
	}()

	return &result
}

func (p *Proxy) start() {
	defer p.rwc.Close()
	r := protocol.NewRequestReader(p.rwc)
	for {
		req, err := r.ReadRequest()
		if err == io.EOF {
			log.Println("eof:", err)
			close(p.closed)
			return
		}
		if err != nil {
			log.Println("proxy:", err)
			close(p.closed)
			return
		}

		resp, err := p.handleRequest(req)
		if err != nil {
			log.Println("request:", err)
			close(p.closed)
			return
		}

		if req.ExtendedSeparator != "" {
			fmt.Fprintln(p.rwc, resp.ExtendedFormat(req.ExtendedSeparator))
		} else {
			fmt.Fprintln(p.rwc, resp.Format())
		}
	}
}

func (p *Proxy) handleRequest(req protocol.Request) (protocol.Response, error) {
	log.Println(">", req.LongFormat())

	if req.Key() == protocol.CommandKey("chk_vfo") {
		log.Println("<", "CHKVFO 0")
		return ChkVfoResponse, nil
	}

	if req.InvalidatesCommand != "" {
		p.cache.Invalidate(req.InvalidatedKey())
	}

	if req.Cacheable {
		resp, ok := p.cache.Get(req.Key())
		if ok {
			log.Println("c", resp.Format())
			return resp, nil
		}
	}

	resp, err := p.trx.Send(context.Background(), req)
	if err != nil {
		return protocol.Response{}, err
	}

	if req.Cacheable {
		p.cache.Put(req.Key(), resp)
	}

	log.Println("<", resp.Format())
	return resp, nil
}

func (p *Proxy) Close() {
	select {
	case <-p.closed:
		return
	default:
		close(p.closed)
	}
}

func (p *Proxy) Wait() {
	<-p.closed
}

type nopCache struct{}

func (c *nopCache) Put(protocol.CommandKey, protocol.Response) {
	// NOP
}

func (c *nopCache) Get(protocol.CommandKey) (protocol.Response, bool) {
	return protocol.Response{}, false
}

func (c *nopCache) Invalidate(protocol.CommandKey) {
	// NOP
}
