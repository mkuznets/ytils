package yfs

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"os"
	"path/filepath"
)

func EnsureDir(parts ...string) (string, error) {
	path := filepath.Join(parts...)

	expandedPath, err := Expand(path)
	if err != nil {
		return "", fmt.Errorf("failed to expand path %s: %w", path, err)
	}

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", absPath, err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %s: %w", absPath, err)
	}

	if err := testWritableDir(absPath); err != nil {
		return "", err
	}

	return absPath, nil
}

func testWritableDir(path string) (err error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("stat path: %s: %w", path, err)
	}

	if !fi.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	f, err := os.CreateTemp(path, ".tmp*")
	if err != nil {
		return fmt.Errorf("write test file: %s: %w", path, err)
	}
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			slog.Warn("remove test file", "err", err)
		}
		if err := f.Close(); err != nil {
			slog.Warn("close test file", "err", err)
		}
	}()

	return nil
}

func Expand(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	if path[0] != '~' {
		return path, nil
	}
	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}
