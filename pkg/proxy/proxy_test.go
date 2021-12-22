package proxy

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/ftl/rigproxy/pkg/test"
)

func TestProxyTransceiverSendReceiveRoundtrip(t *testing.T) {
	proxyBuffer := test.NewBuffer("f\nF 1234\n")
	trxBuffer := test.NewBuffer("get_freq:\n14074000\nRPRT 0\nset_freq: 1234\nRPRT 11\n")

	trx := protocol.NewTransceiver(trxBuffer)
	defer trx.Close()

	proxy := New(proxyBuffer, trx, nil, false)
	defer proxy.Close()
	proxy.Wait()

	trxBuffer.AssertWritten(t, "+\\get_freq\n+\\set_freq 1234\n")
	proxyBuffer.AssertWritten(t, "14074000\nRPRT 11\n")
	proxyBuffer.AssertClosed(t)
}

func TestCommands(t *testing.T) {
	testCases := []struct {
		desc     string
		proxy    string
		trx      string
		expected string
	}{
		{"dump_state", "\\dump_state\n", "dump_state:\n1\n2\n3\nRPRT 0\n", "1\n2\n3\n"},
		{"chk_vfo", "\\chk_vfo\n", "CHKVFO 0\n", "CHKVFO 0\n"},
		{"set_split_vfo short", "S 1 VFOB\n", "set_split_vfo: 1 VFOB\nRPRT 0\n", "RPRT 0\n"},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			trx := protocol.NewTransceiver(test.NewBuffer(tC.trx))
			defer trx.Close()

			proxyBuffer := test.NewBuffer(tC.proxy)
			proxy := New(proxyBuffer, trx, nil, false)
			defer proxy.Close()
			proxy.Wait()

			proxyBuffer.AssertWritten(t, tC.expected)
		})
	}
}

func TestProxyStopsWhenDone(t *testing.T) {
	done := make(chan struct{})
	bytes := make(chan byte)
	r := newChannelReader(bytes)

	proxy := New(r, nil, done, false)

	go func() {
		<-time.After(10 * time.Millisecond)
		close(done)
	}()

	<-done
	proxy.Wait()

}
func TestProxyStopsWhenReceiveFails(t *testing.T) {
	trx := new(mockTransceiver)
	trx.On("Send", mock.Anything, mock.Anything).Once().Return(protocol.Response{}, errors.New("fail"))
	proxy := New(test.NewBuffer("f\n"), trx, nil, false)

	proxy.Wait()
}

func TestProxyInvalidatesCache(t *testing.T) {
	trx := new(mockTransceiver)
	cache := new(mockCache)
	proxy := Proxy{
		trx:   trx,
		cache: cache,
	}

	cache.On("Invalidate", protocol.CommandKey("invalidate_me")).Once()
	trx.On("Send", mock.Anything, mock.Anything).Once().Return(protocol.Response{}, nil)

	proxy.handleRequest(protocol.Request{
		Command: protocol.Command{
			InvalidatesCommand: "invalidate_me",
		},
	})
}

func TestProxyUsesCache(t *testing.T) {
	cache := new(mockCache)
	proxy := Proxy{
		cache: cache,
	}
	cachedCommand := protocol.CommandKey("i_am_cached")
	resp := protocol.Response{
		Data:   []string{"responsedata"},
		Result: "0",
	}

	cache.On("Get", cachedCommand).Once().Return(resp, true)

	actual, err := proxy.handleRequest(protocol.Request{
		Command: protocol.Command{
			Long:      string(cachedCommand),
			Cacheable: true,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, resp, actual)
	cache.AssertExpectations(t)
}

func TestProxyFillsCache(t *testing.T) {
	trx := new(mockTransceiver)
	cache := new(mockCache)
	proxy := Proxy{
		trx:   trx,
		cache: cache,
	}
	cachedCommand := protocol.CommandKey("i_am_cached")
	resp := protocol.Response{
		Data:   []string{"responsedata"},
		Result: "0",
	}

	cache.On("Get", cachedCommand).Once().Return(protocol.Response{}, false)
	trx.On("Send", mock.Anything, mock.Anything).Once().Return(resp, nil)
	cache.On("Put", cachedCommand, resp).Once()

	proxy.handleRequest(protocol.Request{
		Command: protocol.Command{
			Long:      string(cachedCommand),
			Cacheable: true,
		},
	})

	cache.AssertExpectations(t)
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) Put(key protocol.CommandKey, resp protocol.Response) {
	m.Called(key, resp)
}

func (m *mockCache) Get(key protocol.CommandKey) (protocol.Response, bool) {
	args := m.Called(key)
	return args.Get(0).(protocol.Response), args.Bool(1)
}

func (m *mockCache) Invalidate(key protocol.CommandKey) {
	m.Called(key)
}

type mockTransceiver struct {
	mock.Mock
}

func (m *mockTransceiver) Send(ctx context.Context, req protocol.Request) (protocol.Response, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(protocol.Response), args.Error(1)
}

type channelReader struct {
	bytes <-chan byte
	done  chan struct{}
}

func newChannelReader(bytes chan byte) *channelReader {
	return &channelReader{bytes, make(chan struct{})}
}

func (r *channelReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, io.EOF
	}
	select {
	case <-r.done:
		return 0, io.EOF
	case b := <-r.bytes:
		p[0] = b
		return 1, nil
	}
}

func (r *channelReader) Write(p []byte) (int, error) {
	return 0, nil
}

func (r *channelReader) Close() error {
	if _, ok := <-r.done; !ok {
		close(r.done)
	}
	return nil
}
