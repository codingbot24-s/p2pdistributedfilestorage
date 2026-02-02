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

const defaultRoot = "pmcnetwork"

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

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.filename)
}

type StoreOpts struct {
	// folder name of the store containig
	// all the files and folder
	Root string
	PathTransform
}

var DefaultPathTransform = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if len(opts.Root) == 0 {
		opts.Root = defaultRoot
	}
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
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.FullPath())
	return os.Open(fullPathWithRoot)

}

func (s *Store) writeStreams(key string, r io.Reader) error {
	PathKey := s.PathTransform(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.FullPath())
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("written %d bytes: %s\n", n, fullPathWithRoot)
	return nil
}



func (s *Store) Has(key string) bool {
	pathkey := s.PathTransform(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	if err == fs.ErrNotExist {
		return false
	}

	return true
}

func (p *PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (s *Store) Delete(key string) error {
	// get the hash
	pathkey := s.PathTransform(key)
	
	defer func() {
		// here only file
		fmt.Printf("deleting [%s] from disk \n", pathkey.filename)
		
	}()
	// delete the file with root 
	firstPathNamewithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FirstPathName())	
	// TODO: cant remove the main dir  
	os.RemoveAll(firstPathNamewithRoot)

	return nil 

}

func (s *Store) clear()  error {
	return os.RemoveAll(s.Root)
}
