package cmd

import "github.com/rs/zerolog/log"

func Start() {
	log.Info().Msg("starting server")
}

func Stop() {
	log.Log().Msg("stopping server")
}
