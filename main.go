package main

import (
	"log"
	"net"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/ftl/rigproxy/pkg/cache"
	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/ftl/rigproxy/pkg/proxy"
)

var (
	destination = flag.StringP("destination", "d", "localhost:4534", "<host:port> of the destination rigctld server (default: localhost:4534)")
	listen      = flag.StringP("listen", "l", ":4532", "listening address of this proxy (default: :4532)")
	lifetime    = flag.DurationP("lifetime", "L", 200*time.Millisecond, "the lifetime of responses in the cache (default: 200ms)")
)

func main() {
	flag.Parse()

	out, err := net.Dial("tcp", *destination)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	trx := protocol.NewTransceiver(out)

	cache := cache.NewWithLifetime(*lifetime)

	l, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go proxy.NewCached(conn, trx, cache, nil)
	}
}
