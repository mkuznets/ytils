package yhttp

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeJson[T any](r io.Reader) (T, error) {
	var v T
	defer func(src io.Reader) {
		_, _ = io.Copy(io.Discard, src)
	}(r)
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return v, fmt.Errorf("decode JSON: %w", err)
	}
	return v, nil
}

func UnmarshallJson[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return &v, fmt.Errorf("decode JSON: %w", err)
	}
	return &v, nil
}
