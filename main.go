package main

import (
	"context"
	"log"
	"net"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/ftl/rigproxy/pkg/cache"
	"github.com/ftl/rigproxy/pkg/netio"
	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/ftl/rigproxy/pkg/proxy"
)

var (
	destination = flag.StringP("destination", "d", "localhost:4534", "<host:port> of the destination rigctld server (default: localhost:4534)")
	listen      = flag.StringP("listen", "l", ":4532", "listening address of this proxy (default: :4532)")
	lifetime    = flag.DurationP("lifetime", "L", 200*time.Millisecond, "the lifetime of responses in the cache (default: 200ms)")
	timeout     = flag.DurationP("timeout", "t", 10*time.Second, "the timeout for network requests")
	retry       = flag.DurationP("retry", "r", 10*time.Second, "the retry interval")
	trace       = flag.BoolP("trace", "v", false, "trace the communication with the destination")
	test        = flag.BoolP("test", "T", false, "run test code")
)

func main() {
	flag.Parse()

	for {
		if *test {
			runTest()
			return
		}

		loop()
		<-time.After(*retry)
	}
}

func loop() {
	done := make(chan struct{})
	defer func() {
		select {
		case <-done:
		default:
			close(done)
		}
		log.Println("loop done")
	}()

	out, err := net.Dial("tcp", *destination)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	log.Printf("connected to %s", *destination)

	trx := protocol.NewTransceiver(netio.WithTimeout(out, *timeout))
	trx.WhenDone(func() {
		log.Println("transceiver stopped")
		close(done)
	})

	cache := cache.NewWithLifetime(*lifetime)

	l, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		<-done
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go proxy.NewCached(conn, trx, cache, done, *trace)
	}
}

func runTest() {
	out, err := net.Dial("tcp", *destination)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	trx := protocol.NewTransceiver(out)
	trx.WhenDone(func() {
		log.Println("transceiver stopped")
	})

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			request := protocol.Request{Command: protocol.ShortCommand("f")}
			startTime := time.Now()
			response, err := trx.Send(context.Background(), request)
			log.Printf("%v %v", response, time.Now().Sub(startTime))
			if err != nil {
				log.Print("polling frequency failed: ", err)
				return
			}
		}
	}
}
