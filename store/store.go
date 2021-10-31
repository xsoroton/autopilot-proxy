package store

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

// Store interface which can have different implementations: Mem, Redis, DB
type Store interface {
	Set(key string, value []byte) (err error)
	Get(key string) (value []byte, err error)
	Remove(key string) (err error)
}
