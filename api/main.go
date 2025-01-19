package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/mcgovman/wheresmylift/api/cmd"
	"github.com/mcgovman/wheresmylift/api/docs"
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
	docs.SwaggerInfo.Version = Version
}

// @title			WheresMyLift
// @description	Realtime API of the Irish public transit network
// @contact.name	Conor Mc Govern
// @contact.email	wheresmylift(at)mcgov(dot)ie
// @license.name	BSD-3-Clause
// @license.url	https://github.com/mcgovman/wheresmylift/blob/main/LICENSE.md
// @BasePath		/
func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		cmd.Stop()
		os.Exit(0)
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Log().Err(fmt.Errorf("%v", err)).Msg("server panic")
			sigs <- syscall.SIGTERM

			return
		}
	}()

	if len(os.Args) > 1 {
		if os.Args[1] == "-v" {
			fmt.Println(Version)

			return
		} else if _, err := os.Stat(os.Args[1]); errors.Is(err, os.ErrNotExist) {
			log.Error().Err(err).Msg("specified config file does not exist")

			return
		}

		cmd.Start(filepath.Dir(os.Args[1]))
	} else {
		cmd.Start("")
	}
}
