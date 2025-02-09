package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mcgovman/wheresmylift/api/internal/config"
	"github.com/mcgovman/wheresmylift/api/test-utils"
	"github.com/nsf/jsondiff"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func randomAddr() string {
	portNum, _ := rand.Int(rand.Reader, big.NewInt((65000-1024+1)+1024))

	return fmt.Sprintf(":%d", portNum)
}

func validConfig() config.Config {
	return config.Config{
		LogLevel: "debug",
		Timeouts: config.Timeouts{
			Shutdown:   30 * time.Second,
			Startup:    30 * time.Second,
			ReadHeader: 2 * time.Second,
		},
		HTTP: config.HTTP{
			ListenAddress: randomAddr(),
			CORS: config.CORS{
				AllowedOrigins: []string{"*"},
			},
			TrustedProxies: []string{"10.0.0.2"},
		},
	}
}

var assertionStepTimeout time.Duration = 10 * time.Second
var assertionPollInterval time.Duration = 100 * time.Millisecond

func TestStart(t *testing.T) {
	t.Run("cmd will start", func(t *testing.T) {
		Srv = nil
		cfg := validConfig()
		cfgYaml, err := yaml.Marshal(cfg)
		assert.NoError(t, err, "could not marshall config")
		dir, err := os.MkdirTemp("", uuid.New().String())
		assert.NoError(t, err, "could not create temp dir")
		defer os.RemoveAll(dir)
		assert.NoError(t, err, "could not marshal config into yaml")
		cfgPath := path.Join(dir, "api.yml")
		err = os.WriteFile(cfgPath, cfgYaml, 0600)
		assert.NoError(t, err, "could not write config to file")

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		go func() {
			Start(dir)
		}()

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"config":  cfg,
						"message": "got config",
					},
					jsondiff.FullMatch,
				),
				"could not find config log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"level":   "info",
						"message": "starting server",
					},
					jsondiff.FullMatch,
				),
				"could not find server starting log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.Len(t, logSink.Logs, 2, "expected length of logs")
	})

	t.Run("cmd will fail with an invalid config", func(t *testing.T) {
		Srv = nil
		dir, err := os.MkdirTemp("", uuid.New().String())
		defer os.RemoveAll(dir)
		assert.NoError(t, err, "could not create temp dir")
		cfg := validConfig()
		cfg.LogLevel = "some random level"
		cfgYaml, err := yaml.Marshal(cfg)
		assert.NoError(t, err, "could not marshall config")
		err = os.WriteFile(path.Join(dir, "api.yml"), cfgYaml, 0600)
		assert.NoError(t, err, "could not write config to file")

		go func() {
			Start(dir)
		}()

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"config_issues": []string{
							"An invalid log level was specified",
						},
						"message": "configuration issues",
					},
					jsondiff.FullMatch,
				),
				"could not find config issues log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.Len(t, logSink.Logs, 1, "expected length of logs")
	})

	t.Run("cmd will fail with an invalid yaml file", func(t *testing.T) {
		Srv = nil
		dir, err := os.MkdirTemp("", uuid.New().String())
		defer os.RemoveAll(dir)
		assert.NoError(t, err, "could not create temp dir")
		err = os.WriteFile(path.Join(dir, "api.yml"), []byte("a"), 0600)
		assert.NoError(t, err, "could not write config to file")

		go func() {
			Start(dir)
		}()

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"level":   "error",
						"message": "failed to read configuration",
					},
					jsondiff.FullMatch,
				),
				"could not find failed start server log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.Len(t, logSink.Logs, 1, "expected length of logs")
	})

	t.Run("cmd will fail to start the server on an already used port", func(t *testing.T) {
		Srv = nil
		cfg := validConfig()
		cfgYaml, err := yaml.Marshal(cfg)
		assert.NoError(t, err, "could not marshal config into yaml")
		dir, err := os.MkdirTemp("", uuid.New().String())
		defer os.RemoveAll(dir)
		assert.NoError(t, err, "could not create temp dir")
		err = os.WriteFile(path.Join(dir, "api.yml"), cfgYaml, 0600)
		assert.NoError(t, err, "could not write config to file")

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		l, err := net.Listen("tcp", cfg.HTTP.ListenAddress)
		assert.NoError(t, err, "could not create listener")
		defer l.Close()

		go func() {
			Start(dir)
		}()

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"level":   "error",
						"error":   fmt.Sprintf("failed to start HTTP server: listen tcp %s: bind: address already in use", cfg.HTTP.ListenAddress),
						"message": "Failed to start server",
					},
					jsondiff.FullMatch,
				),
				"could not find server start failed log",
			)
		}, assertionStepTimeout, assertionPollInterval)
	})
}

func TestStop(t *testing.T) {
	t.Run("will stop the server", func(t *testing.T) {
		Srv = nil
		cfg := validConfig()
		cfgYaml, err := yaml.Marshal(cfg)
		assert.NoError(t, err, "could not marshal config into yaml")
		dir, err := os.MkdirTemp("", uuid.New().String())
		defer os.RemoveAll(dir)
		assert.NoError(t, err, "could not create temp dir")
		err = os.WriteFile(path.Join(dir, "api.yml"), cfgYaml, 0600)
		assert.NoError(t, err, "could not write config to file")

		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		go func() {
			Start(dir)
		}()

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"level":   "info",
						"message": "starting server",
					},
					jsondiff.FullMatch,
				),
				"could not find server starting log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.Len(t, logSink.Logs, 2, "expected length of logs")

		Stop()
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"message": "stopping server",
					},
					jsondiff.FullMatch,
				),
				"could not find stopping server log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"message": "stopped server successfully",
					},
					jsondiff.FullMatch,
				),
				"could not find stopped server successfully log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.Len(t, logSink.Logs, 4, "expected length of logs")
	})

	t.Run("will stop the server when its not running", func(t *testing.T) {
		Srv = nil
		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		Stop()

		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.True(
				c,
				logSink.ContainsLog(
					map[string]interface{}{
						"message": "stopped server successfully",
					},
					jsondiff.FullMatch,
				),
				"could not find stopped server successfully log",
			)
		}, assertionStepTimeout, assertionPollInterval)
		assert.Len(t, logSink.Logs, 1, "expected length of logs")
	})
}
