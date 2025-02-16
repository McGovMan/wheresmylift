package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Config config.Config
	HTTP   *http.Server
}

func NewServer(config config.Config) *Server {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodOptions,
		},
	})

	r := SetupRouter()
	if config.HTTP.TrustedProxy != "" {
		// The config verifies the IPs are valid
		_ = r.SetTrustedProxies([]string{config.HTTP.TrustedProxy})
	}

	httpSrv := &http.Server{
		Addr:              config.HTTP.ListenAddress,
		Handler:           corsMiddleware.Handler(r),
		ReadHeaderTimeout: 100 * time.Millisecond,
	}

	s := &Server{
		Config: config,
		HTTP:   httpSrv,
	}

	r.GET("", s.RootGet)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("v0/healthcheck", s.V0HealthCheckGet)

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
