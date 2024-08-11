package test

import (
	"gin-demo/app/manager"
	"gin-demo/test/kit"

	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {

	t.Run("health 200", func(t *testing.T) {
		w, _ := kit.HttpGet("/health", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})
}

func TestVersion(t *testing.T) {

	t.Run("version 200", func(t *testing.T) {
		w, _ := kit.HttpGet("/version", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, manager.GetConfig().Env.Version, w.Body.String())
	})
}
