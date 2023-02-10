package yconfig_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mkuznets.com/go/ytils/yconfig"
	"testing"
	"time"
)

type rootCfg struct {
	Foo   string    `config:"cfoo" default:"1"`
	Elems []elemCfg `config:"celems"`
}

func (c rootCfg) Validate() error {
	return nil
}

type elemCfg struct {
	Dur time.Duration `config:"cv" default:"14s"`
}

func TestNewFromMap(t *testing.T) {

	m := map[string]interface{}{
		"celems": []map[string]interface{}{
			{"cv": "1m"},
			{"cv": "10m"},
			{},
		},
	}

	cfg, err := yconfig.NewFromMap[rootCfg](m).Read()

	require.NoError(t, err)
	assert.Equal(t, "1", cfg.Foo)
	assert.Len(t, cfg.Elems, 3)
	assert.Equal(t, time.Minute, cfg.Elems[0].Dur)
	assert.Equal(t, 10*time.Minute, cfg.Elems[1].Dur)
	assert.Equal(t, 14*time.Second, cfg.Elems[2].Dur)
}
