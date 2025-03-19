package gocache

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCacheSet(t *testing.T) {
	// Setup
	c := New()
	c.Set("k2", "value2")

	// Test Case 1: Set new key-value entry
	t.Run("new key-value entry", func(t *testing.T) {
		err := c.Set("k1", "value1")

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}
	})

	// Test Case 2: Set existing key-value entry
	t.Run("existing key entry", func(t *testing.T) {
		err := c.Set("k2", "new value2")

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		if value, _ := c.Get("k2"); value != "new value2" {
			t.Errorf("value - got: %v, want: new value2", value)
		}
	})
}

func TestCacheSetWithTtl(t *testing.T) {
	// Setup
	c := New()
	c.Set("k2", "value2")

	// Test Case 1: Set new key-value entry
	t.Run("new key-value entry", func(t *testing.T) {
		err := c.SetWithTtl("k1", "value1", 100*time.Millisecond)

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		time.Sleep(150 * time.Millisecond)

		if c.Has("k1") {
			t.Errorf("after TTL has passed, key still exists")
		}
	})

	// Test Case 2: Set existing key-value entry
	t.Run("existing key-value entry", func(t *testing.T) {
		err := c.SetWithTtl("k2", "new value2", 100*time.Millisecond)

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		time.Sleep(150 * time.Millisecond)

		if c.Has("k1") {
			t.Errorf("after TTL has passed, key still exists")
		}
	})

	// Test Case 3: Set existing key-value entry with shorter TTL
	t.Run("existing key-value entry with shorter TTL", func(t *testing.T) {
		c.SetWithTtl("k3", "value3", 500*time.Millisecond)

		time.Sleep(100 * time.Millisecond)

		err := c.SetWithTtl("k3", "new value3", 300*time.Millisecond)

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		time.Sleep(250 * time.Millisecond)

		if !c.Has("k3") {
			t.Errorf("before new TTL passes, key was not found")
		}

		time.Sleep(55 * time.Millisecond)

		if c.Has("k3") {
			t.Errorf("new TTL has passed, key still exists")
		}
	})

	// Test Case 4: Set existing key-value entry with longer TTL
	t.Run("existing key-value entry with longer TTL", func(t *testing.T) {
		c.SetWithTtl("k4", "value4", 300*time.Millisecond)

		time.Sleep(100 * time.Millisecond)

		err := c.SetWithTtl("k4", "new value4", 500*time.Millisecond)

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		time.Sleep(205 * time.Millisecond)

		if !c.Has("k4") {
			t.Errorf("before new TTL passes, key was not found")
		}

		time.Sleep(305 * time.Millisecond)

		if c.Has("k4") {
			t.Errorf("new TTL has passed, key still exists")
		}
	})
}

func TestCacheGet(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.SetWithTtl("k2", "value2", 100*time.Millisecond)

	// Test Case 1: Key not found
	t.Run("key not found", func(t *testing.T) {
		value, err := c.Get("nokey")

		if value != nil {
			t.Errorf("value - got: %v, want: nil", value)
		}

		if err == nil || !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("err - got: %v, want: ErrKeyNotFound", err)
		}
	})

	// Test Case 2: Key found
	t.Run("key found", func(t *testing.T) {
		value, err := c.Get("k1")

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		if value == nil || value != "value1" {
			t.Errorf("value - got: %v, want: value1", value)
		}
	})

	// Test Case 3: Key found before and after TTL
	t.Run("key found (before and after TTL)", func(t *testing.T) {
		value, err := c.Get("k2")

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		if value == nil || value != "value2" {
			t.Errorf("value - got: %v, want: value2", value)
		}

		time.Sleep(150 * time.Millisecond)

		value, err = c.Get("k2")

		if value != nil {
			t.Errorf("value - got: %v, want: nil", value)
		}

		if err == nil || !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("err - got: %v, want: ErrKeyNotFound", err)
		}
	})
}

func TestCacheGetAndDelete(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.SetWithTtl("k2", "value2", 100*time.Millisecond)

	// Test Case 1: Key not found
	t.Run("key not found", func(t *testing.T) {
		value, err := c.GetAndDelete("nokey")

		if value != nil {
			t.Errorf("value - got: %v, want: nil", value)
		}

		if err == nil || !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("err - got: %v, want: ErrKeyNotFound", err)
		}
	})

	// Test Case 2: Key found
	t.Run("key found", func(t *testing.T) {
		value, err := c.GetAndDelete("k1")

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		if value == nil || value != "value1" {
			t.Errorf("value - got: %v, want: value1", value)
		}

		if _, ok := c.data["k1"]; ok {
			t.Errorf("value - key still exists")
		}
	})

	// Test Case 3: Key found before and after TTL
	t.Run("key found (before and after TTL)", func(t *testing.T) {
		value, err := c.GetAndDelete("k2")

		if err != nil {
			t.Errorf("err - got: %v, want nil", err)
		}

		if value == nil || value != "value2" {
			t.Errorf("value - got: %v, want: value2", value)
		}

		time.Sleep(150 * time.Millisecond)

		value, err = c.GetAndDelete("k2")

		if value != nil {
			t.Errorf("value - got: %v, want: nil", value)
		}

		if err == nil || !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("err - got: %v, want: ErrKeyNotFound", err)
		}
	})
}

func TestCacheDelete(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")

	// Test Case 1: Delete non-existing key
	t.Run("non-existing key", func(t *testing.T) {
		count := c.Delete("non1")

		if count != 0 {
			t.Errorf("deleted count - got: %d, want 0", count)
		}
	})

	// Test Case 2: Delete existing key
	t.Run("existing key", func(t *testing.T) {
		count := c.Delete("k1")

		if count != 1 {
			t.Errorf("deleted count - got: %d, want 1", count)
		}
	})
}

func TestCacheKeys(t *testing.T) {
	// Setup
	c := New()

	// Test Case 1: No keys
	t.Run("no keys", func(t *testing.T) {
		keys := c.Keys()

		if len(keys) != 0 {
			t.Errorf("keys length - got: %d, want: 0", len(keys))
		}
	})

	// Test Case 2: Existing keys
	t.Run("existing keys", func(t *testing.T) {
		for i := range 3 {
			c.Set(fmt.Sprintf("k%d", i+1), fmt.Sprintf("value%d", i+1))
		}

		keys := c.Keys()

		if len(keys) != 3 {
			t.Errorf("keys length - got: %d, want: 3", len(keys))
		}
	})
}

func TestCacheHas(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")

	// Test Case 1: non-existing key
	t.Run("non-existing key", func(t *testing.T) {
		if c.Has("non1") {
			t.Error("has key non1 - got: true, want: false")
		}
	})

	// Test Case 2: existing key
	t.Run("existing key", func(t *testing.T) {
		if !c.Has("k1") {
			t.Error("has key k1 - got: false, want: true")
		}
	})
}

func TestCacheStats(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.Get("k2")
	c.Get("k1")
	c.SetWithTtl("k2", "value2", 500*time.Millisecond)
	c.Get("k2")

	t.Run("before key TTL", func(t *testing.T) {
		s := c.Stats()

		if s.Hits != 2 {
			t.Errorf("Stats Hits - got: %d, want: 2", s.Hits)
		}

		if s.Misses != 1 {
			t.Errorf("Stats Misses - got: %d, want: 1", s.Misses)
		}

		if s.Keys != 2 {
			t.Errorf("Stats Keys - got: %d, want: 2", s.Keys)
		}

		keySize := SizeOf("k1") + SizeOf("k2")
		valueSize := SizeOf("value1") + SizeOf("value2")

		if s.KeySize != keySize {
			t.Errorf("Stats Key Size - got: %d, want: %d", s.KeySize, keySize)
		}

		if s.ValueSize != valueSize {
			t.Errorf("Stats Value Size - got: %d, want: %d", s.ValueSize, valueSize)
		}
	})

	t.Run("after key TTL", func(t *testing.T) {
		time.Sleep(550 * time.Millisecond)
		s := c.Stats()

		if s.Hits != 2 {
			t.Errorf("Stats Hits - got: %d, want: 2", s.Hits)
		}

		if s.Misses != 1 {
			t.Errorf("Stats Misses - got: %d, want: 1", s.Misses)
		}

		if s.Keys != 1 {
			t.Errorf("Stats Keys - got: %d, want: 1", s.Keys)
		}

		keySize := SizeOf("k1")
		valueSize := SizeOf("value1")

		if s.KeySize != keySize {
			t.Errorf("Stats Key Size - got: %d, want: %d", s.KeySize, keySize)
		}

		if s.ValueSize != valueSize {
			t.Errorf("Stats Value Size - got: %d, want: %d", s.ValueSize, valueSize)
		}
	})

	t.Run("clear stats", func(t *testing.T) {
		c.ClearStats()
		s := c.Stats()

		if s.Hits != 0 {
			t.Errorf("Stats Hits - got: %d, want: 0", s.Hits)
		}

		if s.Misses != 0 {
			t.Errorf("Stats Misses - got: %d, want: 0", s.Misses)
		}

		if s.Keys != 0 {
			t.Errorf("Stats Keys - got: %d, want: 0", s.Keys)
		}

		if s.KeySize != 0 {
			t.Errorf("Stats Key Size - got: %d, want: 0", s.KeySize)
		}

		if s.ValueSize != 0 {
			t.Errorf("Stats Value Size - got: %d, want: 0", s.ValueSize)
		}
	})
}

func TestCacheClear(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.Get("k2")
	c.Get("k1")
	c.SetWithTtl("k2", "value2", 500*time.Millisecond)
	c.Get("k2")

	t.Run("statistics before clearing", func(t *testing.T) {
		s := c.Stats()

		if s.Hits != 2 {
			t.Errorf("Stats Hits - got: %d, want: 2", s.Hits)
		}

		if s.Misses != 1 {
			t.Errorf("Stats Misses - got: %d, want: 1", s.Misses)
		}

		if s.Keys != 2 {
			t.Errorf("Stats Keys - got: %d, want: 2", s.Keys)
		}

		keySize := SizeOf("k1") + SizeOf("k2")
		valueSize := SizeOf("value1") + SizeOf("value2")

		if s.KeySize != keySize {
			t.Errorf("Stats Key Size - got: %d, want: %d", s.KeySize, keySize)
		}

		if s.ValueSize != valueSize {
			t.Errorf("Stats Value Size - got: %d, want: %d", s.ValueSize, valueSize)
		}
	})

	t.Run("clear cache", func(t *testing.T) {
		c.Clear()
		s := c.Stats()

		if s.Hits != 0 {
			t.Errorf("Stats Hits - got: %d, want: 0", s.Hits)
		}

		if s.Misses != 0 {
			t.Errorf("Stats Misses - got: %d, want: 0", s.Misses)
		}

		if s.Keys != 0 {
			t.Errorf("Stats Keys - got: %d, want: 0", s.Keys)
		}

		if s.KeySize != 0 {
			t.Errorf("Stats Key Size - got: %d, want: 0", s.KeySize)
		}

		if s.ValueSize != 0 {
			t.Errorf("Stats Value Size - got: %d, want: 0", s.ValueSize)
		}

		if c.Has("k1") {
			t.Error("Has \"k1\" - got: true, want: false")
		}

		if c.Has("k2") {
			t.Error("Has \"k2\" - got: true, want: false")
		}
	})
}
