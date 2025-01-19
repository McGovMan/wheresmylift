package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/mcgovman/wheresmylift/api/test-utils"
	"github.com/nsf/jsondiff"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var RemoteAddr = "10.0.0.2"
var ClientAddr = "10.0.0.3"

func TestStart(t *testing.T) {
	t.Run("should use the ip in the X-Real-IP when the trusted proxy is set", func(t *testing.T) {
		cfg := config.Config{
			HTTP: config.HTTP{
				TrustedProxies: []string{RemoteAddr},
			},
		}
		srv := NewServer(cfg)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = fmt.Sprintf("%s:1234", RemoteAddr)
		r.Header.Set("X-Real-IP", ClientAddr)

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		srv.HTTP.Handler.ServeHTTP(w, r)
		assert.True(
			t,
			logSink.ContainsLog(
				map[string]interface{}{
					"level":   "info",
					"method":  "GET",
					"path":    "/",
					"status":  307,
					"ip":      ClientAddr,
					"message": "request_info",
				},
				jsondiff.SupersetMatch,
			),
		)
	})

	t.Run("the server will succeed running on a port", func(t *testing.T) {
		portNum, err := rand.Int(rand.Reader, big.NewInt((65000-1024+1)+1024))
		assert.NoError(t, err)
		port := fmt.Sprintf(":%d", portNum)
		cfg := config.Config{
			HTTP: config.HTTP{
				ListenAddress: port,
			},
		}
		srv := NewServer(cfg)

		defer func() {
			srv.HTTP.Close()
		}()

		go func() {
			startErr := srv.Start()
			assert.NoError(t, startErr, "server should be able to start")
		}()

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%s", port), time.Second)
			assert.NoError(c, err, "should be able to check if port is being used")
			assert.NotNil(c, conn, "port should be active")
			if conn != nil {
				conn.Close()
			}
		}, time.Second, 50*time.Millisecond)
	})

	t.Run("the server will fail running on a port", func(t *testing.T) {
		cfg := config.Config{
			HTTP: config.HTTP{
				ListenAddress: ":80",
			},
		}
		srv := NewServer(cfg)

		defer func() {
			srv.HTTP.Close()
		}()

		err := srv.Start()
		assert.Error(t, err, "server should not be able to start")
	})

	t.Run("the server responds on an endpoint", func(t *testing.T) {
		srv := NewServer(config.Config{})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/v0/healthcheck", nil)

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		srv.HTTP.Handler.ServeHTTP(w, r)
		assert.Equal(t, 204, w.Result().StatusCode, "expect response status to be 204")
		assert.True(
			t,
			logSink.ContainsLog(
				map[string]interface{}{
					"level":   "info",
					"method":  "GET",
					"path":    "/v0/healthcheck",
					"status":  204,
					"message": "request_info",
				},
				jsondiff.SupersetMatch,
			),
		)
		assert.Len(t, logSink.Logs, 1, "expect only one log entry")
	})
}

func TestStop(t *testing.T) {
	t.Run("can shut down a server successfully", func(t *testing.T) {
		portNum, err := rand.Int(rand.Reader, big.NewInt((65000-1024+1)+1024))
		assert.NoError(t, err)
		port := fmt.Sprintf(":%d", portNum)
		cfg := config.Config{
			HTTP: config.HTTP{
				ListenAddress: port,
			},
		}
		srv := NewServer(cfg)

		defer func() {
			srv.HTTP.Close()
		}()

		go func() {
			srv.HTTP.ListenAndServe()
		}()

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%s", port), time.Second)
			assert.NoError(c, err, "should be able to check if port is being used")
			assert.NotNil(c, conn, "port should be active")
			if conn != nil {
				conn.Close()
			}
		}, time.Second, 50*time.Millisecond)

		srv.Stop(context.Background())

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%s", port), time.Second)
			assert.Error(c, err, "should get connection refused")
			assert.Nil(c, conn, "port should be inactive")
			if conn != nil {
				conn.Close()
			}
		}, time.Second, 50*time.Millisecond)
	})
}
