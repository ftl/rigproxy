package proxy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/ftl/rigproxy/pkg/test"
)

func TestProxyTransceiverSendReceiveRoundtrip(t *testing.T) {
	clientBuffer := test.NewBuffer("f\nF 1234\n")
	trxBuffer := test.NewBuffer("14074000\nRPRT 0\nRPRT 11\n")

	trx := protocol.NewTransceiver(trxBuffer)
	defer trx.Close()

	proxy := New(clientBuffer, trx)
	defer proxy.Close()
	proxy.Wait()

	trxBuffer.AssertWritten(t, "\\get_freq\n\\set_freq 1234\n")
	clientBuffer.AssertWritten(t, "14074000\nRPRT 0\nRPRT 11\n")
	clientBuffer.AssertClosed(t)
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

	actual := proxy.handleRequest(protocol.Request{
		Command: protocol.Command{
			Long:      string(cachedCommand),
			Cacheable: true,
		},
	})

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
