package hw04lrucache

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
		require.Equal(t, 100, val.([]interface{})[1])

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val.([]interface{})[1])

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val.([]interface{})[1])

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		_, ok := c.Get("a")
		require.True(t, ok)
		_, ok = c.Get("b")
		require.True(t, ok)
		_, ok = c.Get("c")
		require.True(t, ok)

		c.Set("d", 4)

		_, ok = c.Get("a")
		require.False(t, ok)

		_, ok = c.Get("b")
		require.True(t, ok)
		_, ok = c.Get("c")
		require.True(t, ok)
		_, ok = c.Get("d")
		require.True(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		c.Set("a", 11)
		v, ok := c.Get("a")
		require.Equal(t, v.([]interface{})[1], 11)
		require.True(t, ok)

		c.Set("c", 33)
		v, ok = c.Get("c")
		require.Equal(t, v.([]interface{})[1], 33)
		require.True(t, ok)

		c.Get("a")
		c.Get("c")
		c.Set("d", 4)

		_, ok = c.Get("b")
		require.False(t, ok)
		_, ok = c.Get("a")
		require.True(t, ok)
		_, ok = c.Get("c")
		require.True(t, ok)
		_, ok = c.Get("d")
		require.True(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
