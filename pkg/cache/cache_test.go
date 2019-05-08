package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestEmptyCache(t *testing.T) {
	cache := New()

	_, ok := cache.Get(theCommand)

	assert.False(t, ok)
}

func TestPutGetRoundtrip(t *testing.T) {
	cache := New()
	resp := protocol.Response{
		Data:   []string{"response_data"},
		Result: "0",
	}

	cache.Put(theCommand, resp)
	actual, ok := cache.Get(theCommand)

	assert.True(t, ok)
	assert.Equal(t, resp, actual)
}

func TestConcurrentAccess(t *testing.T) {
	cache := New()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			resp := protocol.Response{
				Data: []string{fmt.Sprint(i)},
			}
			cache.Put(theCommand, resp)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			resp, ok := cache.Get(theCommand)
			if ok && resp.Data[0] == "999" {
				return
			}
		}
	}()

	wg.Wait()
}

func TestInvalidate(t *testing.T) {
	cache := New()
	resp := protocol.Response{
		Data:   []string{"response_data"},
		Result: "0",
	}

	cache.Put(theCommand, resp)
	cache.Invalidate(theCommand)
	_, ok := cache.Get(theCommand)

	assert.False(t, ok)
}

func TestLifetime(t *testing.T) {
	cache := NewWithLifetime(10 * time.Millisecond)
	resp := protocol.Response{
		Data:   []string{"response_data"},
		Result: "0",
	}

	cache.Put(theCommand, resp)
	actual, ok := cache.Get(theCommand)

	assert.True(t, ok)
	assert.Equal(t, resp, actual)

	time.Sleep(10 * time.Millisecond)
	actual, ok = cache.Get(theCommand)

	assert.False(t, ok)
}

const theCommand = protocol.CommandKey("the_command")
