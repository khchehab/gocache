package gocache

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

const keyPoolSize = 1024

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
	c.SetWithTtl("k2", "value2", 300*time.Millisecond)

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

	// Test Case 3: Delete existing key with TTL
	t.Run("existing key with ttl", func(t *testing.T) {
		count := c.Delete("k2")

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

func TestCacheClear(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.Get("k2")
	c.Get("k1")
	c.SetWithTtl("k2", "value2", 500*time.Millisecond)
	c.Get("k2")

	t.Run("clear cache", func(t *testing.T) {
		c.Clear()

		if c.Has("k1") {
			t.Error("Has \"k1\" - got: true, want: false")
		}

		if c.Has("k2") {
			t.Error("Has \"k2\" - got: true, want: false")
		}
	})
}

func TestCacheGetTtl(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.SetWithTtl("k2", "value3", 500*time.Millisecond)

	if ttl := c.GetTtl("k1"); ttl != 0 {
		t.Errorf("GetTtl of \"k1\" - got: %v, want: 0", ttl)
	}

	if ttl := c.GetTtl("k2"); ttl != 500*time.Millisecond {
		t.Errorf("GetTtl of \"k2\" - got: %v, want: 500ms", ttl)
	}

	if ttl := c.GetTtl("nonexistent"); ttl != -1 {
		t.Errorf("GetTtl of \"nonexistent\" - got: %v, want: -1", ttl)
	}

	time.Sleep(550 * time.Millisecond)

	if ttl := c.GetTtl("k2"); ttl != -1 {
		t.Errorf("GetTtl of \"k2\" after TTL - got: %v, want: -1", ttl)
	}
}

func TestCacheChangeTtl(t *testing.T) {
	// Setup
	c := New()
	c.Set("k1", "value1")
	c.SetWithTtl("k2", "value2", 500*time.Millisecond)
	c.Set("k3", "value3")
	c.SetWithTtl("k4", "value4", 700*time.Millisecond)

	// Test Case 1: Change TTL to -1
	t.Run("ChangeTtl to -1", func(t *testing.T) {
		if !c.ChangeTtl("k3", -1) {
			t.Errorf("ChangeTtl of \"k3\" to -1 - got: false, want: true")
		}
	})

	// Test Case 2: Change TTL from 0 to 200ms
	t.Run("ChangeTtl from 0 to 200ms", func(t *testing.T) {
		if !c.ChangeTtl("k1", 200*time.Millisecond) {
			t.Errorf("ChangeTtl of \"k1\" to 200ms - got: false, want: true")
		}

		if ttl := c.GetTtl("k1"); ttl != 200*time.Millisecond {
			t.Errorf("GetTtl of \"k1\" - got: %v, want: 200ms", ttl)
		}

		time.Sleep(250 * time.Millisecond)

		if c.Has("k1") {
			t.Errorf("\"k1\" still exists after new TTL")
		}
	})

	// Test Case 3: Change TTL from 500ms to 300ms
	t.Run("ChangeTtl from 500ms to 300ms", func(t *testing.T) {
		if !c.ChangeTtl("k2", 300*time.Millisecond) {
			t.Errorf("ChangeTtl of \"k2\" to 300ms - got: false, want: true")
		}

		if ttl := c.GetTtl("k2"); ttl != 300*time.Millisecond {
			t.Errorf("GetTtl of \"k2\" - got: %v, want: 300ms", ttl)
		}

		time.Sleep(350 * time.Millisecond)

		if c.Has("k2") {
			t.Errorf("\"k2\" still exists after new TTL")
		}
	})

	// Test Case 4: Change TTL from 700ms to 900ms
	t.Run("ChangeTtl from 700ms to 900ms", func(t *testing.T) {
		if !c.ChangeTtl("k4", 900*time.Millisecond) {
			t.Errorf("ChangeTtl of \"k4\" to 700ms - got: false, want: true")
		}

		if ttl := c.GetTtl("k4"); ttl != 900*time.Millisecond {
			t.Errorf("GetTtl of \"k4\" - got: %v, want: 700ms", ttl)
		}

		time.Sleep(550 * time.Millisecond)

		if !c.Has("k4") {
			t.Errorf("\"k4\" does not exist before new TTL")
		}

		time.Sleep(400 * time.Millisecond)

		if c.Has("k4") {
			t.Errorf("\"k4\" still exists after new TTL")
		}
	})

	// Test Case 5: Change TTL of non existent
	t.Run("ChangeTtl of non-existent", func(t *testing.T) {
		if c.ChangeTtl("nonexistent", 200*time.Millisecond) {
			t.Errorf("ChangeTtl of \"nonexistent\" to 200ms - got: true, want: false")
		}
	})
}

func BenchmarkCacheSet(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	values := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		values[i] = fmt.Sprintf("value%d", i)
	}

	b.ResetTimer()

	for i := range b.N {
		c.Set(keys[i%keyPoolSize], values[i%keyPoolSize])
	}
}

func BenchmarkCacheSetWithTtl(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	values := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		values[i] = fmt.Sprintf("value%d", i)
	}

	b.ResetTimer()

	for i := range b.N {
		c.SetWithTtl(keys[i%keyPoolSize], values[i%keyPoolSize], 100*time.Millisecond)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		c.Set(keys[i], fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		c.Get(keys[i%keyPoolSize])
	}
}

func BenchmarkCacheGetAndDelete(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		c.Set(keys[i], fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		c.GetAndDelete(keys[i%keyPoolSize])
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		c.Set(keys[i], fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		c.Delete(keys[i%keyPoolSize])
	}
}

func BenchmarkCacheKeys(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
	}

	b.ResetTimer()

	for range b.N {
		c.Keys()
	}
}

func BenchmarkCacheHas(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		c.Set(keys[i], fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		c.Has(keys[i%keyPoolSize])
	}
}

func BenchmarkCacheClear(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	values := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		values[i] = fmt.Sprintf("value%d", i)
	}

	b.ResetTimer()

	for range b.N {
		b.StopTimer()

		for i := range keyPoolSize {
			c.Set(keys[i], values[i])
		}

		b.StartTimer()

		c.Clear()
	}
}

func BenchmarkCacheGetTtl(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		c.Set(keys[i], fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		c.GetTtl(keys[i%keyPoolSize])
	}
}

func BenchmarkCacheChangeTtl(b *testing.B) {
	c := New()

	keys := make([]string, keyPoolSize)
	for i := range keyPoolSize {
		keys[i] = strconv.Itoa(i)
		c.Set(keys[i], fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		c.ChangeTtl(keys[i%keyPoolSize], 100*time.Millisecond)
	}
}
