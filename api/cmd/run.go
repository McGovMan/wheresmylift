package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/mcgovman/wheresmylift/api/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var srv *server.Server

func reload(sigs chan os.Signal) {
	if srv != nil {
		stop()
		srv = nil
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Msg("failed to read configuration")
		sigs <- syscall.SIGTERM

		return
	}

	var cfg config.Config
	// This will only error if you don't pass it a pointer
	_ = viper.Unmarshal(&cfg)
	issues := cfg.Verify()
	if len(issues) != 0 {
		log.Log().Strs("config_issues", issues).Msg("configuration issues")
		sigs <- syscall.SIGTERM

		return
	}

	log.Log().Any("config", cfg).Msg("got config")

	logLevel := cfg.GetZeroLogLevel()
	zerolog.SetGlobalLevel(logLevel)

	srv = server.NewServer(cfg)

	log.Info().Msg("starting server")
	if err := srv.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to start server")
		sigs <- syscall.SIGTERM

		return
	}
}

func stop() {
	if srv != nil {
		log.Log().Msg("stopping server")

		ctx, cancel := context.WithTimeout(context.Background(), srv.Config.Timeouts.Shutdown)
		defer cancel()
		srv.Stop(ctx)
	}

	log.Log().Msg("stopped server successfully")
}

// @title          	WheresMyLift
// @description     Realtime API of the Irish public transit network
// @contact.name   	Conor Mc Govern
// @contact.email  	wheresmylift(at)mcgov(dot)ie
// @license.name 	BSD-3-Clause
// @license.url   	https://github.com/mcgovman/wheresmylift/blob/main/LICENSE.md
// @BasePath  		/
func Run(sigs chan os.Signal, configDir string) {
	// Recover from panic
	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(fmt.Errorf("%v", err)).Msg("server panic")
			sigs <- syscall.SIGTERM

			return
		}
	}()

	viper.SetConfigType("yml")
	viper.SetConfigName("api")
	viper.AddConfigPath("/run/config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configDir)

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Log().Str("file", e.Name).Msg("config changed - reloading")
		reload(sigs)
	})
	viper.WatchConfig()
	reload(sigs)

	<-sigs
	stop()
}
