package yconfig

import (
	"fmt"
	"github.com/kkyr/fig"
)

type config interface {
	Validate() error
}

type Reader[T config] struct {
	filename   string
	lookupDirs []string
}

func New[T config](filename string) *Reader[T] {
	return &Reader[T]{
		filename:   filename,
		lookupDirs: []string{"."},
	}
}

func (r *Reader[T]) WithLookupDir(dir string) *Reader[T] {
	r.lookupDirs = append(r.lookupDirs, dir)
	return r
}

func (r *Reader[T]) Read() (T, error) {
	var cfg T
	err := fig.Load(
		&cfg,
		fig.File(r.filename),
		fig.Dirs(r.lookupDirs...),
		fig.Tag("config"),
	)
	if err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
