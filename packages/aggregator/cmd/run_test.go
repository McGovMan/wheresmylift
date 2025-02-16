package cmd

import (
	"testing"
	"time"

	"github.com/mcgovman/wheresmylift/lib/go-test-utils"
	"github.com/nsf/jsondiff"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var assertionStepTimeout time.Duration = 10 * time.Second
var assertionPollInterval time.Duration = 100 * time.Millisecond

func TestStart(t *testing.T) {
	t.Run("cmd will start", func(t *testing.T) {
		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		go func() {
			Start()
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
		assert.Len(t, logSink.Logs, 1, "expected length of logs")
	})
}

func TestStop(t *testing.T) {
	t.Run("will stop the server", func(t *testing.T) {
		logSink := test.LogSink{}
		log.Logger = zerolog.New(&logSink)

		go func() {
			Start()
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
		assert.Len(t, logSink.Logs, 1, "expected length of logs")

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
		assert.Len(t, logSink.Logs, 2, "expected length of logs")
	})
}
