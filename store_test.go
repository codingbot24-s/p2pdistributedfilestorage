package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransform: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.clear(); err != nil {
		t.Error(err)
	}

}

func TestPathTransformFunc(t *testing.T) {
	key := "mykey"
	pathName := CASPathTransformFunc(key)
	fmt.Println(pathName)
}



func TestStore(t *testing.T) {

	s := newStore()
	defer tearDown(t, s)
	key := "mykey"
	data := []byte("hello world")
	if err := s.writeStreams(key, bytes.NewReader(data)); err != nil {
		t.Fatal(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have a key %s", key)
	}

	r, err := s.read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("got %q, want %q", string(b), string(data))
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
