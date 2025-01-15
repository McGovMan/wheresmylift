package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/rs/cors"
)

type Server struct {
	Config config.Config
	HTTP   *http.Server
}

func NewServer(config config.Config) *Server {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: config.HTTP.CORS.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodOptions,
		},
	})

	r := SetupRouter()
	if config.HTTP.TrustedProxies != nil {
		// The config verifies the IPs are valid
		_ = r.SetTrustedProxies(config.HTTP.TrustedProxies)
	}

	httpSrv := &http.Server{
		Addr:              config.HTTP.ListenAddress,
		Handler:           corsMiddleware.Handler(r),
		ReadHeaderTimeout: config.Timeouts.ReadHeader,
	}

	s := &Server{
		Config: config,
		HTTP:   httpSrv,
	}

	r.GET("", func(c *gin.Context) {
		c.Status(204)
	})

	return s
}

func (s *Server) Start() error {
	if err := s.HTTP.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	// Ignoring the errors here as it is difficult to test that the shutdown will fail
	// To counter this we attempt to shutdown gracefully, and if we can't, we forcefully do so
	_ = s.HTTP.Shutdown(ctx)
	_ = s.HTTP.Close()
}
