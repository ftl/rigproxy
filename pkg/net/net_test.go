package net

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutgoingConnectionClosedRemotely(t *testing.T) {
	remoteListening := make(chan struct{})
	var remoteAddr string
	go func() {
		l, err := net.Listen("tcp", "")
		require.NoError(t, err)
		remoteAddr = l.Addr().String()
		close(remoteListening)
		conn, err := l.Accept()
		require.NoError(t, err)
		conn.Close()
	}()

	<-remoteListening
	conn, err := net.Dial("tcp", remoteAddr)
	require.NoError(t, err)
	p := make([]byte, 10)
	_, err = conn.Read(p)
	assert.Equal(t, io.EOF, err)
}

func TestStopAcceptingNewConnections(t *testing.T) {
	l, err := net.Listen("tcp", "")
	require.NoError(t, err)
	go func() {
		<-time.After(10 * time.Millisecond)
		l.Close()
	}()
	_, err = l.Accept()
	assert.Error(t, err)
}
