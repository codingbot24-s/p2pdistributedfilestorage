package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mykey"
	pathName := CASPathTransformFunc(key)
	fmt.Println(pathName)
}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransform: CASPathTransformFunc,
	}
	store := NewStore(opts)
	key := "mykey"
	data := []byte("hello world")
	if err := store.writeStreams(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := store.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransform: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "mykey"

	data := []byte("hello world")
	if err := s.writeStreams(key, bytes.NewReader(data)); err != nil {
		t.Fatal(err)
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

	s.Delete(key)
}
