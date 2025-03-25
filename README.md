# gocache

[![Go Reference](https://pkg.go.dev/badge/github.com/khchehab/gocache.svg)](https://pkg.go.dev/github.com/khchehab/gocache)
[![Test](https://github.com/khchehab/gocache/actions/workflows/go.yml/badge.svg)](https://github.com/khchehab/gocache/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/khchehab/gocache)](https://goreportcard.com/report/github.com/khchehab/gocache)
![GitHub Release](https://img.shields.io/github/v/release/khchehab/gocache)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

gocache is an in-memory cache in Go, it aims to provide a fast and efficient store to quickly read/write data.

## Installation

```shell
go get github.com/khchehab/gocache
```

## Usage

```go
package main

import (
    "log"
    "github.com/khchehab/gocache"
)

func main() {
    // Create a new instance of Cache
    cache := gocache.New()

    // Set a key-value pair
    if err := cache.Set("prevId", 1); err != nil {
        log.Fatalf("error setting key 'prevId' in cache: %v", err)
    }

    if err := cache.SetWithTtl("paused", true, 5 * time.Second); err != nil {
        log.Fatalf("error setting key 'paused' in cache: %v", err)
    }

    // Get value from cache
    value, err := cache.Get("paused")
    if err != nil {
        log.Fatalf("error getting key 'paused' from cache: %v", err)
    }

    log.Printf("paused from cache: %v\n", value)

    // Delete from cache
    cache.Delete("prevId")
}
```

## Customization

You can customize the cache that you create with options.
- Global TTL: You can set the global cache's time-to-live for all the entries (can be overriden for a single pair). A value of `0` means unlimited and will remain until manually removed. Any other value will be based on the duration set.

```go
func main() {
    // all entries will be valid for 5 seconds.
    cache := gocache.New(gocache.WithStdTtl(5 * time.Second))
}
```

- Delete on expiration: A flag to indicate whether an entry should be automatically deleted from the store after it has expired. `true` means the entry will be deleted from the store after it's time has passed. `false` means the entry will remain but will be flagged as expired and will be treated as non-existent.

```go
func main() {
    // expired entries will stay in store but flagged as expired
    cache := gocache.New(gocache.WithDeleteOnExpire(false))
}
```

- Maximum number of keys: A number that indicates a maximum number of entries the cache can hold. A value of `-1` means unlimited keys. Any other number will enforce the number of keys that can be stored.

```go
func main() {
    // set a maximum of 10 keys in the cache
    cache := gocache.New(gocache.WithMaxKeys(10))
}
```

## Functionalities to Add

Below are some functionalities that I plan to add:
- [ ] Add SyncCache for a concurrent-safe caching.
- [ ] Add `Multiple*` function for operations that deal with multiple entries at the same time.
- [ ] Add a `ForEach` function that loops over the entries and calls a function on each entry.

## License

This project is licensed under the [MIT license](LICENSE).
