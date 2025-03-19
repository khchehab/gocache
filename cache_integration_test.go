package gocache

import (
	"testing"
	"time"
)

func TestCacheIntegration(t *testing.T) {
	c := New(
		WithStdTtl(500*time.Millisecond),
		WithMaxKeys(3),
		WithDeleteOnExpire(false),
	)

	// Set values
	if err := c.Set("k1", "value1"); err != nil {
		t.Errorf("Set k1: err - got: %v, want: nil", err)
	}

	if err := c.SetWithTtl("k2", "value2", 300*time.Millisecond); err != nil {
		t.Errorf("Set k2: err - got: %v, want: nil", err)
	}

	if err := c.SetWithTtl("k3", "value3", 0); err != nil {
		t.Errorf("Set k3: err - got: %v, want: nil", err)
	}

	if err := c.Set("k4", "value4"); err == nil || err != ErrCacheFull {
		t.Error("Set k4: err - got: nil, want: ErrCacheFull")
	}

	// Get values
	value, err := c.Get("k1")

	if value == nil || value != "value1" {
		t.Errorf("Get k1: val - got: %v, want: value1", value)
	}

	if err != nil {
		t.Errorf("Get k1: val - got: %v, want: nil", err)
	}

	value, err = c.Get("k2")

	if value == nil || value != "value2" {
		t.Errorf("Get k2: val - got: %v, want: value2", value)
	}

	if err != nil {
		t.Errorf("Get k2: val - got: %v, want: nil", err)
	}

	value, err = c.Get("k3")

	if value == nil || value != "value3" {
		t.Errorf("Get k3: val - got: %v, want: value3", value)
	}

	if err != nil {
		t.Errorf("Get k3: val - got: %v, want: nil", err)
	}

	value, err = c.GetAndDelete("nonexistent")

	if value != nil {
		t.Errorf("Get nonexistent: val - got: %v, want: nil", value)
	}

	if err == nil || err != ErrKeyNotFound {
		t.Errorf("Get nonexistent: val - got: %v, want: ErrKeyNotFound", err)
	}

	// Check stats
	s := c.Stats()

	if s.Hits != 3 {
		t.Errorf("Stats Hits - got: %v, want: 3", s.Hits)
	}

	if s.Misses != 1 {
		t.Errorf("Stats Misses - got: %v, want: 1", s.Misses)
	}

	if s.Keys != 3 {
		t.Errorf("Stats Keys - got: %v, want: 3", s.Keys)
	}

	keySize := SizeOf("k1") + SizeOf("k2") + SizeOf("k3")
	valueSize := SizeOf("value1") + SizeOf("value2") + SizeOf("value3")

	if s.KeySize != keySize {
		t.Errorf("Stats KeySize - got: %v, want: %v", s.KeySize, keySize)
	}

	if s.ValueSize != valueSize {
		t.Errorf("Stats ValueSize - got: %v, want: %v", s.ValueSize, valueSize)
	}

	// Sleep for some time, before all entries expire
	time.Sleep(350 * time.Millisecond)

	// Clear stats
	c.ClearStats()
	s = c.Stats()

	if s.Hits != 0 {
		t.Errorf("Stats Hits - got: %v, want: 0", s.Hits)
	}

	if s.Misses != 0 {
		t.Errorf("Stats Misses - got: %v, want: 0", s.Misses)
	}

	if s.Keys != 0 {
		t.Errorf("Stats Keys - got: %v, want: 0", s.Keys)
	}

	if s.KeySize != 0 {
		t.Errorf("Stats KeySize - got: %v, want: 0", s.KeySize)
	}

	if s.ValueSize != 0 {
		t.Errorf("Stats ValueSize - got: %v, want: 0", s.ValueSize)
	}

	// Get values
	value, err = c.Get("k1")

	if value == nil || value != "value1" {
		t.Errorf("Get k1: val - got: %v, want: value1", value)
	}

	if err != nil {
		t.Errorf("Get k1: val - got: %v, want: nil", err)
	}

	value, err = c.Get("k2")

	if value != nil {
		t.Errorf("Get k2: val - got: %v, want: nil", value)
	}

	if err == nil || err != ErrKeyNotFound {
		t.Errorf("Get k2: val - got: %v, want: ErrKeyNotFound", err)
	}

	value, err = c.Get("k3")

	if value == nil || value != "value3" {
		t.Errorf("Get k3: val - got: %v, want: value3", value)
	}

	if err != nil {
		t.Errorf("Get k3: val - got: %v, want: nil", err)
	}

	// Sleep for more time for all entries with TTL to expire
	time.Sleep(200 * time.Millisecond)

	// Get values
	value, err = c.Get("k1")

	if value != nil {
		t.Errorf("Get k1: val - got: %v, want: nil", value)
	}

	if err == nil || err != ErrKeyNotFound {
		t.Errorf("Get k1: val - got: %v, want: ErrKeyNotFound", err)
	}

	value, err = c.Get("k2")

	if value != nil {
		t.Errorf("Get k2: val - got: %v, want: nil", value)
	}

	if err == nil || err != ErrKeyNotFound {
		t.Errorf("Get k2: val - got: %v, want: ErrKeyNotFound", err)
	}

	value, err = c.Get("k3")

	if value == nil || value != "value3" {
		t.Errorf("Get k3: val - got: %v, want: value3", value)
	}

	if err != nil {
		t.Errorf("Get k3: val - got: %v, want: nil", err)
	}

	// Check and change TTL of entry
	if ttl := c.GetTtl("k3"); ttl != 0 {
		t.Errorf("GetTtl k3 - got: %v, want: 0", ttl)
	}

	if !c.ChangeTtl("k3", 600*time.Millisecond) {
		t.Errorf("ChangeTtl k3 - got: false, want: true")
	}

	// Get the keys
	keys := c.Keys()

	if len(keys) != 3 {
		t.Errorf("Keys length - got: %v; want: 3", len(keys))
	}

	// Check key existence
	if c.Has("k1") {
		t.Error("Has k1 - got: true, want: false")
	}

	if !c.Has("k3") {
		t.Error("Has k3 - got: false, want: true")
	}

	// Check stats again
	s = c.Stats()

	if s.Hits != 3 {
		t.Errorf("Stats Hits - got: %v, want: 3", s.Hits)
	}

	if s.Misses != 3 {
		t.Errorf("Stats Misses - got: %v, want: 3", s.Misses)
	}

	if s.Keys != 0 {
		t.Errorf("Stats Keys - got: %v, want: 0", s.Keys)
	}

	if s.KeySize != 0 {
		t.Errorf("Stats KeySize - got: %v, want: 0", s.KeySize)
	}

	if s.ValueSize != 0 {
		t.Errorf("Stats ValueSize - got: %v, want: 0", s.ValueSize)
	}

	// Delete an entry
	if count := c.Delete("nonexistent"); count != 0 {
		t.Errorf("Delete nonexistent - got: %v, want: 0", count)
	}

	if count := c.Delete("k3"); count != 1 {
		t.Errorf("Delete k3 - got: %v, want: 1", count)
	}

	if c.Has("k3") {
		t.Errorf("Has k3 - got: true, want: false")
	}

	// Set an item then clear whole cache
	c.Set("k4", "value4")

	if len(c.data) == 0 {
		t.Error("cache data length - got: 0, want: different than 0")
	}

	c.Clear()

	if len(c.data) != 0 {
		t.Errorf("cache data length - got: %v, want: 0", len(c.data))
	}

	// Set an entry with TTL and delete it
	c.Set("k5", "value5")

	time.Sleep(550 * time.Millisecond)

	value, err = c.GetAndDelete("k5")

	if value != nil {
		t.Errorf("Get k5: val - got: %v, want: nil", value)
	}

	if err == nil || err != ErrKeyNotFound {
		t.Errorf("Get k5: val - got: %v, want: ErrKeyNotFound", err)
	}
}
