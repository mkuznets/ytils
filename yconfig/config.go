package yconfig

import (
	"fmt"
	"mkuznets.com/go/ytils/yconfig/internal/fig"
)

type config interface {
	Validate() error
}

type Reader[T config] struct {
	dirs []string
	opts []fig.Option
}

func New[T config](filename string) *Reader[T] {
	return &Reader[T]{
		opts: []fig.Option{
			fig.File(filename),
		},
		dirs: []string{"."},
	}
}

func NewFromMap[T config](m map[string]interface{}) *Reader[T] {
	return &Reader[T]{
		opts: []fig.Option{
			fig.UseMap(m),
			fig.IgnoreFile(),
		},
		dirs: []string{"."},
	}
}

func (r *Reader[T]) WithLookupDir(dir string) *Reader[T] {
	r.dirs = append(r.dirs, dir)
	return r
}

func (r *Reader[T]) Read() (T, error) {
	var cfg T
	opts := append(r.opts,
		fig.Dirs(r.dirs...),
		fig.Tag("config"),
	)

	err := fig.Load(&cfg, opts...)
	if err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
