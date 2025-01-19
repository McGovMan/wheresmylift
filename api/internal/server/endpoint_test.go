package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRootGet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/", s.RootGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "expected status 307 from endpoint")
		assert.Equal(t, "/docs/index.html", w.HeaderMap.Get("Location"), "unexcepted or missing redirect")
	})
}

func TestV0HealthCheckGet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v0/healthcheck", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v0/healthcheck", s.V0HealthCheckGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code, "expected status 204 from endpoint")
	})
}
