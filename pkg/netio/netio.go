package netio

import (
	"io"
	"net"
	"strings"
	"time"
)

// WithTimeout returns a ReadWriteCloser that applies the given timeout to read and write operations.
func WithTimeout(conn net.Conn, timeout time.Duration) io.ReadWriteCloser {
	return &readWriteCloserWithTimeout{
		conn:         conn,
		readTimeout:  timeout,
		writeTimeout: timeout,
	}
}

// IsTimeout indicates if the given error is a timeout error.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasSuffix(err.Error(), "i/o timeout")
}

// WithReadTimeout returns a ReadWriteCloser that applies the given timeout only to read operations.
func WithReadTimeout(conn net.Conn, timeout time.Duration) io.ReadWriteCloser {
	return &readWriteCloserWithTimeout{
		conn:        conn,
		readTimeout: timeout,
	}
}

// WithWriteTimeout returns a ReadWriteCloser that applies the given timeout only to write operations.
func WithWriteTimeout(conn net.Conn, timeout time.Duration) io.ReadWriteCloser {
	return &readWriteCloserWithTimeout{
		conn:         conn,
		writeTimeout: timeout,
	}
}

type readWriteCloserWithTimeout struct {
	conn         net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (c *readWriteCloserWithTimeout) Read(p []byte) (int, error) {
	if c.readTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	return c.conn.Read(p)
}

func (c *readWriteCloserWithTimeout) Write(p []byte) (int, error) {
	if c.writeTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	return c.conn.Write(p)
}

func (c *readWriteCloserWithTimeout) Close() error {
	return c.conn.Close()
}
