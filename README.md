# Carbon Cache
A wrapper around [BadgerDB](https://github.com/dgraph-io/badger) providing a simple API.

**NOTE** This package is provided "as is" with no guarantee. Use it at your own risk and always test it yourself before using it in a production environment. If you find any issues, please [create a new issue](https://github.com/twiny/carbon/issues/new).

## Install
`go get https://github.com/twiny/carbon`

## API
```go
Set(key string, val []byte, ttl time.Duration) error
Get(key string) ([]byte, error)
Del(key string) error
ForEach(prefix string, fn func(key string, val []byte) error) error
```

## Usage
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/twiny/carbon"
)

func main() {
	cache, err := carbon.NewCache("./tmp")
	if err != nil {
		log.Println(err)
		return
	}
	defer cache.Close()

	// set
	if err := cache.Set("foo", []byte("bar"), 10*time.Minute); err != nil {
		log.Println(err)
		return
	}

	// get
	val, err := cache.Get("foo")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(val))

	// range
	if err := cache.ForEach("", func(key string, val []byte) error {
		fmt.Println(key, string(val))
		return nil
	}); err != nil {
		log.Println(err)
		return
	}

	// delete
	if err := cache.Del("foo"); err != nil {
		log.Println(err)
		return
	}
}
```