package carbon

import (
	"testing"
	"time"
)

// TODO (twiny) add tests

func TestSet(t *testing.T) {
	cache, err := NewCache("./_test")
	if err != nil {
		t.FailNow()
		return
	}
	defer cache.Close()

	// test cases
	cases := []struct {
		name string
		item struct {
			key string
			val []byte
		}
		got  interface{}
		want interface{}
	}{
		{
			name: "key/val is empty",
			item: struct {
				key string
				val []byte
			}{key: "", val: []byte("")},
		},
		{
			name: "not found",
			item: struct {
				key string
				val []byte
			}{key: "", val: []byte("")},
		},
	}

	// run
	for _, c := range cases {
		if err := cache.Set(c.item.key, c.item.val, 1*time.Minute); err != nil {
			t.Errorf(err.Error())
			return
		}
	}
}

func TestGet(t *testing.T) {}

func TestDel(t *testing.T) {}

func BenchmarkSet(b *testing.B) {
	cache, err := NewCache("./_test")
	if err != nil {
		b.Fatal(err)
		return
	}
	defer cache.Close()
	//
	for n := 0; n < b.N; n++ {
		// strconv.Itoa(n)
		if err := cache.Set("hello", []byte("world"), 10*time.Minute); err != nil {
			b.Fatal(err)
			return
		}
	}
}

func BenchmarkGet(b *testing.B) {
	cache, err := NewCache("./_test")
	if err != nil {
		b.Fatal(err)
		return
	}
	defer cache.Close()
	//
	for n := 0; n < b.N; n++ {
		// strconv.Itoa(n) // "hello"
		_, err := cache.Get("hello")
		if err != nil {
			b.Fatal(err)
			return
		}
	}
}
