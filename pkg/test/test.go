package test

import (
	"bytes"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ReadWriteCloseBuffer struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	closingLock *sync.Mutex
	closed      bool
}

func NewBuffer(read string) *ReadWriteCloseBuffer {
	return &ReadWriteCloseBuffer{
		readBuffer:  bytes.NewBufferString(read),
		writeBuffer: bytes.NewBufferString(""),
		closingLock: new(sync.Mutex),
	}
}

func (b *ReadWriteCloseBuffer) Read(p []byte) (n int, err error) {
	return b.readBuffer.Read(p)
}

func (b *ReadWriteCloseBuffer) Write(p []byte) (n int, err error) {
	return b.writeBuffer.Write(p)
}

func (b *ReadWriteCloseBuffer) Close() error {
	b.closingLock.Lock()
	defer b.closingLock.Unlock()
	b.closed = true
	return nil
}

func (b *ReadWriteCloseBuffer) LocalAddr() net.Addr {
	return nil
}

func (b *ReadWriteCloseBuffer) RemoteAddr() net.Addr {
	return nil
}

func (b *ReadWriteCloseBuffer) SetDeadline(time.Time) error {
	return nil
}

func (b *ReadWriteCloseBuffer) SetReadDeadline(time.Time) error {
	return nil
}

func (b *ReadWriteCloseBuffer) SetWriteDeadline(time.Time) error {
	return nil
}

func (b *ReadWriteCloseBuffer) AssertWritten(t *testing.T, expected string) {
	assert.Equal(t, expected, b.writeBuffer.String())
}

func (b *ReadWriteCloseBuffer) AssertClosed(t *testing.T) {
	b.closingLock.Lock()
	defer b.closingLock.Unlock()
	assert.True(t, b.closed)
}
