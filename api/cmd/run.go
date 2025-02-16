package cmd

import (
	"context"

	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/mcgovman/wheresmylift/api/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var Srv *server.Server

func Start() {
	viper.SetEnvPrefix("WML")
	viper.AutomaticEnv()

	logLevel := viper.GetString("LOG_LEVEL")
	httpListenAddr := viper.GetString("HTTP_LISTEN_ADDRESS")
	httpTrustedProxy := viper.GetString("HTTP_TRUSTED_PROXY")

	cfg := config.Config{
		LogLevel: logLevel,
		HTTP: config.HTTP{
			ListenAddress: httpListenAddr,
			TrustedProxy:  httpTrustedProxy,
		},
	}

	issues := cfg.Verify()
	if len(issues) != 0 {
		log.Log().Strs("config_issues", issues).Msg("configuration issues")

		return
	}

	log.Log().Any("config", cfg).Msg("got config")

	zerologLevel := cfg.GetZeroLogLevel()
	zerolog.SetGlobalLevel(zerologLevel)

	Srv = server.NewServer(cfg)

	log.Info().Msg("starting server")
	if err := Srv.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to start server")

		return
	}
}

func Stop() {
	if Srv != nil {
		log.Log().Msg("stopping server")
		Srv.Stop(context.Background())
	}

	log.Log().Msg("stopped server successfully")
}
