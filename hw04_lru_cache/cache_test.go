package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300) // ccc, bbb, aaa
		_ = c.Set("ddd", 400) // ddd, ccc, bbb
		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		c = NewCache(3)
		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300) // ccc, bbb, aaa
		_, _ = c.Get("aaa") // aaa, ccc, bbb
		_ = c.Set("ddd", 400) // ddd, aaa, ccc
		// bbb should purged
		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
		// aaa still in cache
		val, ok = c.Get("aaa") // aaa, ddd, ccc
		require.True(t, ok)
		require.Equal(t, val, 100)
		// bbb not in cache
		wasInCache := c.Set("bbb", 200)
		require.False(t, wasInCache)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300)
		c.Clear()
		val, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})	
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
