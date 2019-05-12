# rigproxy

The rigproxy connects to a [Hamlib](https://github.com/Hamlib/Hamlib) `rigctld` server and accepts incoming connections from clients that speak the Hamlib net (model #2) protocol. Incoming requests are forwarded to the destination server, responses of reading requests are cached for a configurable amount of time.

The main purpose of rigproxy is to reduce to load on the destination server and rig by reducing the amount of reading requests if there run multiple clients concurrently.

## Usage

`rigproxy` provides a CLI with the following flags:

* --destination -d <host:port> # the address of the destination `rigctld` server
* --listen -l <if:port> # the listening interface and port, `if` may be empty to bind to all available network interfaces
* --lifetime -L <duration> # the duration that responses to reading requests are cached

For example:

```
rigproxy -d localhost:4534 -l :4532 -L 200ms
```


## Disclaimer
I develop this tool for myself and just for fun in my free time. If you find it useful, I'm happy to hear about that. If you have trouble using it, you have all the source code to fix the problem yourself (although pull requests are welcome).

## License
This tool is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/) 2019