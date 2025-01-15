package main

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/mcgovman/wheresmylift/api/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version string = "dev"

func init() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.
		With().Caller().Logger().
		With().Str("version", Version).Logger()
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	var configFilePath string

	if len(os.Args) > 1 {
		configFilePath = os.Args[1]
	}

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) && configFilePath != "" {
		log.Fatal().Err(err).Msg("specified config file does not exist")
	}

	cmd.Run(sigs, filepath.Dir(configFilePath))
	<-sigs
}
