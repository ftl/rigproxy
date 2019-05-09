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

var ServerError = protocol.Response{Result: "500"}

func New(rwc io.ReadWriteCloser, trx Transceiver) *Proxy {
	return NewCached(rwc, trx, new(nopCache))
}

func NewCached(rwc io.ReadWriteCloser, trx Transceiver, cache Cache) *Proxy {
	result := Proxy{
		rwc:    rwc,
		trx:    trx,
		cache:  new(nopCache),
		closed: make(chan struct{}),
	}

	go result.start()

	return &result
}

func (p *Proxy) start() {
	defer p.rwc.Close()
	r := protocol.NewRequestReader(p.rwc)
	for {
		select {
		case <-p.closed:
			return
		default:
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

			resp := p.handleRequest(req)

			if req.ExtendedSeparator != "" {
				fmt.Fprintln(p.rwc, resp.ExtendedFormat(req.ExtendedSeparator))
			} else {
				fmt.Fprintln(p.rwc, resp.Format())
			}
		}
	}
}

func (p *Proxy) handleRequest(req protocol.Request) protocol.Response {
	log.Println(">", req.LongFormat())
	if req.InvalidatesCommand != "" {
		p.cache.Invalidate(req.InvalidatedKey())
	}

	if req.Cacheable {
		resp, ok := p.cache.Get(req.Key())
		if ok {
			log.Println("c", resp.Format())
			return resp
		}
	}

	resp, err := p.trx.Send(context.Background(), req)
	if err != nil {
		log.Println("request:", err)
		return ServerError
	}

	if req.Cacheable {
		p.cache.Put(req.Key(), resp)
	}

	log.Println("<", resp.Format())
	return resp
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
