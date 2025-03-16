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
	for i := range 10 {
		if i%2 == 0 {
			c.SetWithTtl(fmt.Sprintf("k%d", i+1), fmt.Sprintf("value%d", i+1), 100*time.Millisecond)
		} else {
			c.Set(fmt.Sprintf("k%d", i+1), fmt.Sprintf("value%d", i+1))
		}
	}

	// Test Case 1: Delete single non-existing key
	t.Run("non-existing key", func(t *testing.T) {
		count := c.Delete("non1")

		if count != 0 {
			t.Errorf("deleted count - got: %d, want 0", count)
		}
	})

	// Test Case 2: Delete all non-existing keys
	t.Run("all non-existing keys", func(t *testing.T) {
		count := c.Delete("non1", "non2", "non3")

		if count != 0 {
			t.Errorf("deleted count - got: %d, want 0", count)
		}
	})

	// Test Case 3: Delete single existing key
	t.Run("existing key", func(t *testing.T) {
		count := c.Delete("k1")

		if count != 1 {
			t.Errorf("deleted count - got: %d, want 1", count)
		}
	})

	// Test Case 4: Delete existing and non-existing keys
	t.Run("existing and non-existing keys", func(t *testing.T) {
		count := c.Delete("k1", "k2", "k3", "non1")

		if count != 2 {
			t.Errorf("deleted count - got: %d, want 2", count)
		}
	})

	// Test Case 5: Delete all existing keys
	t.Run("all existing keys", func(t *testing.T) {
		count := c.Delete("k4", "k5", "k6", "k7", "k8", "k9", "k10")

		if count != 7 {
			t.Errorf("deleted count - got: %d, want 7", count)
		}
	})

	// Test Case 6: No keys passed
	t.Run("no keys passed", func(t *testing.T) {
		count := c.Delete()

		if count != 0 {
			t.Errorf("deleted count - got: %d, want 0", count)
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
