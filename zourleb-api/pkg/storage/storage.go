// Package storage abstracts object storage (local disk or S3-compatible) behind
// a single interface. The API stores only URLs.
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Storage saves bytes and returns a public URL.
type Storage interface {
	Save(ctx context.Context, key string, r io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}

// LocalStorage writes to a directory and serves via a configured base URL.
type LocalStorage struct {
	root      string
	publicURL string
}

func NewLocal(root, publicURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &LocalStorage{root: root, publicURL: strings.TrimRight(publicURL, "/")}, nil
}

func (s *LocalStorage) Save(_ context.Context, key string, r io.Reader, _ string) (string, error) {
	dst := filepath.Join(s.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return s.PublicURL(key), nil
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(s.root, filepath.FromSlash(key)))
}

func (s *LocalStorage) PublicURL(key string) string {
	return s.publicURL + "/" + strings.TrimLeft(key, "/")
}

// Key builds a collision-resistant object key under a prefix.
func Key(prefix, filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("%s/%d%s", strings.Trim(prefix, "/"), time.Now().UnixNano(), ext)
}
