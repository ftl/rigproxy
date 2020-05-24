# rigproxy

The rigproxy connects to a [Hamlib](https://github.com/Hamlib/Hamlib) `rigctld` server and accepts incoming connections from clients that speak the Hamlib net (model #2) protocol. Incoming requests are forwarded to the destination server, responses of reading requests are cached for a configurable amount of time.

The main purpose of rigproxy is to reduce the load on the destination server and rig by reducing the amount of reading requests if there run multiple clients concurrently.

In addition, rigproxy also contains a client library that allows to access `rigctld` servers through a native Go programming interface:

```go
conn, err := client.Open("localhost:4532")
if err != nil {
	log.Fatal(err)
}
defer conn.Close()

frequency, err := conn.Frequency(context.Background())
if err != nil {
	log.Fatal(err)
}
log.Printf("current frequency: %.0fHz", frequency)
```

See [godoc](https://godoc.org/github.com/ftl/rigproxy/pkg/client) for more information.

## Usage

`rigproxy` provides a CLI with the following flags:

* --destination -d <host:port> # the address of the destination `rigctld` server
* --listen -l <if:port> # the listening interface and port, `if` may be empty to bind to all available network interfaces
* --lifetime -L <duration> # the duration that responses to reading requests are cached

For example:

```
rigproxy -d localhost:4534 -l :4532 -L 200ms
```

## Development

To use your local copy of rigproxy in other projects, put the following into the go.mod file of your project:

```
replace github.com/ftl/rigproxy => <path_to_your_local_copy>
```

## Disclaimer
I develop this tool for myself and just for fun in my free time. If you find it useful, I'm happy to hear about that. If you have trouble using it, you have all the source code to fix the problem yourself (although pull requests are welcome).

## License
This tool is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)