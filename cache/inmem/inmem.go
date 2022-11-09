package inmem

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/realnighthawk/bucky/cache"
)

type inmem struct {
	handler *gocache.Cache
}

type Options struct {
	Expiration      time.Duration
	CleanupInterval time.Duration
}

func New(opts Options) (cache.Handler, error) {
	h := gocache.New(opts.Expiration, opts.CleanupInterval)
	return &inmem{
		handler: h,
	}, nil
}

func (h *inmem) Get(key string) (interface{}, error) {
	res, ok := h.handler.Get(key)
	if !ok {
		return res, ErrKeyNotExist
	}
	return res, nil
}

func (h *inmem) Set(key string, value interface{}, exp ...time.Duration) error {
	if len(exp) == 0 {
		h.handler.Set(key, value, gocache.NoExpiration)
	} else {
		h.handler.Set(key, value, exp[0])
	}
	return nil
}
