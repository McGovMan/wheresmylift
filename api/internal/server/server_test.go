package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var RemoteAddr = "10.0.0.2"
var ClientAddr = "10.0.0.3"

func TestNewServer(t *testing.T) {
	t.Run("should run the server on the specified port", func(t *testing.T) {
		portNum, err := rand.Int(rand.Reader, big.NewInt((65000-1024+1)+1024))
		assert.NoError(t, err)
		port := fmt.Sprintf(":%d", portNum)
		cfg := config.Config{
			HTTP: config.HTTP{
				ListenAddress: port,
			},
		}
		srv := NewServer(cfg)

		assert.Equal(t, port, srv.HTTP.Addr, "configured address should be on http server")
	})

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

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		srv.HTTP.Handler.ServeHTTP(w, r)
		var logResult map[string]string
		_ = json.Unmarshal(buf.Bytes(), &logResult)

		assert.Equal(t, ClientAddr, logResult["ip"])
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

		go func() {
			time.Sleep(1 * time.Second)
			srv.Stop(context.Background())
		}()

		startErr := srv.Start()
		assert.NoError(t, startErr, "server should be able to start")
	})

	t.Run("the server will fail running on a port", func(t *testing.T) {
		cfg := config.Config{
			HTTP: config.HTTP{
				ListenAddress: ":80",
			},
		}
		srv := NewServer(cfg)

		go func() {
			time.Sleep(1 * time.Second)
			srv.Stop(context.Background())
		}()

		err := srv.Start()
		assert.Error(t, err, "server should not be able to start")
	})
}
