package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashstr := hex.EncodeToString(hash[:])
	blockSize := 5
	strlen := len(hashstr) / blockSize
	paths := make([]string, strlen)

	for i := range strlen {
		paths[i] = hashstr[i*blockSize : (i+1)*blockSize]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		filename: hashstr,
	}
}

type PathTransform func(string) PathKey

type PathKey struct {
	PathName string
	filename string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.filename)
}

type StoreOpts struct {
	PathTransform
}

var DefaultPathTransform = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {

	PathKey := s.PathTransform(key)
	return os.Open(PathKey.FullPath())

}

func (s *Store) writeStreams(key string, r io.Reader) error {
	PathKey := s.PathTransform(key)

	if err := os.MkdirAll(PathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	fullPath := PathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("written %d bytes: %s\n", n, fullPath)
	return nil
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransform(key)

	_, err := os.Stat(pathkey.FullPath())
	if err == fs.ErrNotExist {
		return false
	}

	return true
}
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (s *Store) Delete(key string) error {
	// get the hash
	path := s.PathTransform(key)
	// delete with os.Remove return os.Remove
	defer func() {
		// here only file
		fmt.Printf("deleting [%s] from disk \n", path.filename)
	}()
	return os.RemoveAll(path.FirstPathName())
}
