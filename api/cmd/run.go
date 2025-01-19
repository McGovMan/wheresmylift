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

func Start(configDir string) {
	viper.SetConfigType("yml")
	viper.SetConfigName("api")
	viper.AddConfigPath("/run")
	viper.AddConfigPath(".")
	if configDir != "" {
		viper.AddConfigPath(configDir)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Msg("failed to read configuration")

		return
	}

	var cfg config.Config
	// This will only error if you don't pass it a pointer
	_ = viper.Unmarshal(&cfg)
	issues := cfg.Verify()
	if len(issues) != 0 {
		log.Log().Strs("config_issues", issues).Msg("configuration issues")

		return
	}

	log.Log().Any("config", cfg).Msg("got config")

	logLevel := cfg.GetZeroLogLevel()
	zerolog.SetGlobalLevel(logLevel)

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

		ctx, cancel := context.WithTimeout(context.Background(), Srv.Config.Timeouts.Shutdown)
		defer cancel()
		Srv.Stop(ctx)
	}

	log.Log().Msg("stopped server successfully")
}
