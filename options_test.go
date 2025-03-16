package gocache

import (
	"testing"
	"time"
)

var maxKeysTestCases = []struct {
	label    string
	opt      OptFunc
	expected int
}{
	{"without opts", nil, -1},
	{"negative max keys", WithMaxKeys(-10), -1},
	{"zero max keys", WithMaxKeys(0), 0},
	{"positive max keys", WithMaxKeys(10), 10},
}

var stdTtlTestCases = []struct {
	label    string
	opt      OptFunc
	expected time.Duration
}{
	{"without opts", nil, 0},
	{"zero ttl", WithStdTtl(0), 0},
	{"positive ttl", WithStdTtl(30), 30},
}

var deleteOnExpireTestCases = []struct {
	label    string
	opt      OptFunc
	expected bool
}{
	{"without opts", nil, true},
	{"with delete on expire true", WithDeleteOnExpire(true), true},
	{"with delete on expire false", WithDeleteOnExpire(false), false},
}

func TestMaxKeysOpts(t *testing.T) {
	for _, tc := range maxKeysTestCases {
		t.Run(tc.label, func(t *testing.T) {
			var c *Cache
			if tc.opt != nil {
				c = New(tc.opt)
			} else {
				c = New()
			}

			if c.maxKeys != tc.expected {
				t.Errorf("maxKeys - got: %v, want: %v", c.maxKeys, tc.expected)
			}
		})
	}
}

func TestStdTtlOpts(t *testing.T) {
	for _, tc := range stdTtlTestCases {
		t.Run(tc.label, func(t *testing.T) {
			var c *Cache
			if tc.opt != nil {
				c = New(tc.opt)
			} else {
				c = New()
			}

			if c.stdTtl != tc.expected {
				t.Errorf("stdTtl - got: %v, want: %v", c.stdTtl, tc.expected)
			}
		})
	}
}

func TestDeleteOnExpireOpts(t *testing.T) {
	for _, tc := range deleteOnExpireTestCases {
		t.Run(tc.label, func(t *testing.T) {
			var c *Cache
			if tc.opt != nil {
				c = New(tc.opt)
			} else {
				c = New()
			}

			if c.deleteOnExpire != tc.expected {
				t.Errorf("deleteOnExpire - got: %v, want: %v", c.deleteOnExpire, tc.expected)
			}
		})
	}
}
