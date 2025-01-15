package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	t.Run("check setup router", func(t *testing.T) {
		r := SetupRouter()
		assert.Len(t, r.Handlers, 2, "should include 2 middlewares from engine")
		assert.Equal(t, r.BasePath(), "/", "base path should be /")
	})

	t.Run("recovery with valid recovered interface", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, r)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/recovery", nil)
		assert.NoError(t, err, "could not create http request")

		r.GET("/recovery", func(c *gin.Context) {
			panic("Oh no :(")
		})

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		r.ServeHTTP(w, req)
		var logResult map[string]string
		err = json.Unmarshal(buf.Bytes(), &logResult)
		assert.NoError(t, err, "could not unmarshal logging result to interface")
		assert.Equal(t, "Oh no :(", logResult["error"], "expected panic message to be logged")
	})

	t.Run("should generate a new context id if not included", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, r)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")
		r.GET("/", func(c *gin.Context) {})

		r.ServeHTTP(w, req)
		_, contextIdErr := uuid.Parse(ctx.Writer.Header().Get("context-id"))
		assert.NoError(t, contextIdErr, "context id should be a valid uuid")
	})

	t.Run("should use existing context id if included", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, r)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		req.Header.Set("context-id", "388240b8-2826-435a-9556-03f6da0c8894")
		assert.NoError(t, err, "could not create http request")
		r.GET("/", func(c *gin.Context) {})

		r.ServeHTTP(w, req)
		assert.Equal(t, "388240b8-2826-435a-9556-03f6da0c8894", ctx.Writer.Header().Get("context-id"), "existing context id should be included")
	})

	t.Run("should log the method, path, response status, latency_ns, request_info, and context id", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, r)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
		assert.NoError(t, err, "could not create http request")
		r.GET("/test", func(c *gin.Context) {
			c.Status(202)
		})

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		r.ServeHTTP(w, req)
		var logResult map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &logResult)
		assert.NoError(t, err, "could not unmarshal logging result to interface")
		assert.Equal(t, float64(202), logResult["status"], "expected status to be logged")
		assert.Equal(t, "/test", logResult["path"], "expected path to be logged")
		assert.Equal(t, http.MethodGet, logResult["method"], "expected method to be logged")
		assert.Regexp(t, regexp.MustCompile(`[\d]+`), logResult["latency_ns"], "expected latency_ns to be logged")
		_, contextIdErr := uuid.Parse(ctx.Writer.Header().Get("context-id"))
		assert.NoError(t, contextIdErr, "context id should be a valid uuid")
		assert.Equal(t, "request_info", logResult["message"], "should log log type")
	})

	t.Run("should use client IP if proxy is not used", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, r)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")
		r.GET("/", func(c *gin.Context) {})
		err = r.SetTrustedProxies(nil)
		assert.NoError(t, err, "could not set trusted proxies to nil")
		req.RemoteAddr = "10.0.0.2:1234"
		req.Header.Set("X-Real-IP", "10.0.0.3")

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		r.ServeHTTP(w, req)
		var logResult map[string]interface{}
		_ = json.Unmarshal(buf.Bytes(), &logResult)
		assert.Equal(t, "10.0.0.2", logResult["ip"], "should log client ip")
	})

	t.Run("should use X-Real-Ip if proxy is used", func(t *testing.T) {
		r := SetupRouter()
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, r)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")
		r.GET("/", func(c *gin.Context) {})
		err = r.SetTrustedProxies([]string{"10.0.0.2"})
		assert.NoError(t, err, "could not set trusted proxies")
		req.RemoteAddr = "10.0.0.2:1234"
		req.Header.Set("X-Real-IP", "10.0.0.3")

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		r.ServeHTTP(w, req)
		var logResult map[string]interface{}
		_ = json.Unmarshal(buf.Bytes(), &logResult)
		assert.Equal(t, "10.0.0.3", logResult["ip"], "should log X-Real-Ip ip")
	})
}
