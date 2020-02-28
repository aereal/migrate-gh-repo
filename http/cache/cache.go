package cache

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

var ErrCacheMissed = errors.New("missed")

type Cache interface {
	Get(r *http.Request) (*http.Response, error)
	Set(r *http.Request, resp *http.Response) error
}

type nopCache int

var NopCache Cache = nopCache(0)

func (c nopCache) Get(r *http.Request) (*http.Response, error) {
	return nil, ErrCacheMissed
}

func (c nopCache) Set(r *http.Request, resp *http.Response) error {
	return nil
}

type FileCache struct {
	Root string
}

func (c *FileCache) Get(r *http.Request) (*http.Response, error) {
	key, err := canonicalCacheKey(r)
	if err != nil {
		return nil, fmt.Errorf("cannot create cache key: %w", err)
	}
	path := filepath.Join(c.Root, key)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache file %s: %w", path, err)
	}
	return http.ReadResponse(bufio.NewReader(f), r)
}

func (c *FileCache) Set(r *http.Request, resp *http.Response) error {
	key, err := canonicalCacheKey(r)
	if err != nil {
		return fmt.Errorf("cannot create cache key: %w", err)
	}
	path := filepath.Join(c.Root, key)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	content, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return fmt.Errorf("failed to dump response: %w", err)
	}
	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("failed to write response to file %s: %w", path, err)
	}
	return nil
}

func canonicalCacheKey(r *http.Request) (string, error) {
	bytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(bytes)
	return fmt.Sprintf("%x", h), nil
}
