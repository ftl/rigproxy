package main

import (
	"log"
	"net"

	"github.com/ftl/rigproxy/pkg/cache"
	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/ftl/rigproxy/pkg/proxy"
)

func main() {
	out, err := net.Dial("tcp", "localhost:4534")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	trx := protocol.NewTransceiver(out)
	cache := cache.New()

	l, err := net.Listen("tcp", ":4532")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go proxy.NewCached(conn, trx, cache)
	}
}
