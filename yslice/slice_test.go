package yslice_test

import (
	"github.com/stretchr/testify/assert"
	"mkuznets.com/go/ytils/yslice"
	"testing"
)

func TestMapByKey(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		values := make([]string, 0)
		result := yslice.MapByKey(values, func(value string) string {
			return value[:1]
		})
		assert.Equal(t, result, map[string][]string{})
	})

	t.Run("duplicate", func(t *testing.T) {
		values := []string{"a1", "a2", "b2", "c3"}

		result := yslice.MapByKey(values, func(value string) string {
			return value[:1]
		})
		assert.Equal(t, result["a"], []string{"a1", "a2"})
		assert.Equal(t, result["b"], []string{"b2"})
		assert.Equal(t, result["c"], []string{"c3"})
	})
}
