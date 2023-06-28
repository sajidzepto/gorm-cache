package inmemory

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"time"
)

type Store struct {
	store *cache.Cache
}

func New() *Store {
	inMemCache := cache.New(10*time.Hour, 10*time.Minute)
	return &Store{store: inMemCache}
}

// NewWithDb
// @param tx
func NewWithDb(tx *cache.Cache) *Store {
	return &Store{store: tx}
}

// Set
// @param ctx
// @param key
// @param value
// @param ttl
// @date 2022-07-02 08:12:11
func (r *Store) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	r.store.Set(key, value, ttl)
	return nil
}

// Get
// @param ctx
// @param key
func (r *Store) Get(ctx context.Context, key string) ([]byte, error) {
	get, found := r.store.Get(key)
	if found {
		if g, ok := get.([]uint8); ok {
			return []byte(g), nil
		} else {
			return nil, errors.New("cannot deserialize key")
		}
	}
	return nil, errors.New("key not found")
}

// RemoveFromTag
// @param ctx
// @param tag
func (r *Store) RemoveFromTag(ctx context.Context, tag string) error {
	r.store.Set(tag, map[string]bool{}, -1)
	return nil
}

// SaveTagKey
// @param ctx
// @param tag
// @param key
func (r *Store) SaveTagKey(ctx context.Context, tag, key string) error {
	get, found := r.store.Get(tag)
	if !found {
		r.store.Set(tag, map[string]bool{key: true}, -1)
		return nil
	} else {
		if m, ok := get.(map[string]bool); ok {
			m[key] = true
			r.store.Set(tag, m, -1)
			return nil
		} else {
			return errors.New("cannot add in tag set")
		}
	}
}
