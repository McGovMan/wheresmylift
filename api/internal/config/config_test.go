package config

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetZeroLogLevel(t *testing.T) {
	t.Run("check all levels can be converted", func(t *testing.T) {
		c := &Config{}

		c.LogLevel = "trace"
		assert.Equal(t, zerolog.TraceLevel, c.GetZeroLogLevel(), "expected to return trace zerolog ENUM")

		c.LogLevel = "disabled"
		assert.Equal(t, zerolog.Disabled, c.GetZeroLogLevel(), "expected to return disabled zerolog ENUM")

		c.LogLevel = "panic"
		assert.Equal(t, zerolog.PanicLevel, c.GetZeroLogLevel(), "expected to return panic zerolog ENUM")

		c.LogLevel = "fatal"
		assert.Equal(t, zerolog.FatalLevel, c.GetZeroLogLevel(), "expected to return fatal zerolog ENUM")

		c.LogLevel = "error"
		assert.Equal(t, zerolog.ErrorLevel, c.GetZeroLogLevel(), "expected to return error zerolog ENUM")

		c.LogLevel = "warn"
		assert.Equal(t, zerolog.WarnLevel, c.GetZeroLogLevel(), "expected to return warn zerolog ENUM")

		c.LogLevel = "info"
		assert.Equal(t, zerolog.InfoLevel, c.GetZeroLogLevel(), "expected to return info zerolog ENUM")

		c.LogLevel = "debug"
		assert.Equal(t, zerolog.DebugLevel, c.GetZeroLogLevel(), "expected to return debug zerolog ENUM")

		c.LogLevel = "dummy"
		assert.Equal(t, zerolog.NoLevel, c.GetZeroLogLevel(), "expected to return nolevel zerolog ENUM")
	})
}

var validConfig Config = Config{
	LogLevel: "debug",
	HTTP: HTTP{
		ListenAddress: ":8080",
		TrustedProxy:  "127.0.0.1",
	},
}

type Run struct {
	name        string
	beforeWork  func()
	verifyFunc  func() []string
	issue       string
	expectIssue bool
}

func (r *Run) verifyIssuesAndError(t *testing.T) {
	r.beforeWork()
	issues := r.verifyFunc()
	if r.expectIssue {
		assert.Contains(t, issues, r.issue, "expected an issue but was not found")
	} else {
		assert.NotContains(t, issues, r.issue, "expected no issue but was found")
	}
}

func TestHTTPVerify(t *testing.T) {
	var testConfig Config

	runs := []Run{
		// HTTP: Listen Address
		{
			name:        "expect no HTTP listen address issue",
			beforeWork:  func() {},
			expectIssue: false,
		},
		{
			name: "expect HTTP listen address issue when given nothing",
			beforeWork: func() {
				testConfig.HTTP.ListenAddress = ""
			},
			issue:       "HTTP listen address is not valid",
			expectIssue: true,
		},
		{
			name: "expect HTTP listen address issue when given nothing",
			beforeWork: func() {
				testConfig.HTTP.ListenAddress = "...abc"
			},
			issue:       "HTTP listen address is not valid",
			expectIssue: true,
		},
		// HTTP: Trusted proxies
		{
			name: "expect no trusted proxies issue with valid IP",
			beforeWork: func() {
				testConfig.HTTP.TrustedProxy = "127.0.0.1"
			},
			expectIssue: false,
		},
		{
			name: "expect trusted proxies issue with invalid IP",
			beforeWork: func() {
				testConfig.HTTP.TrustedProxy = "a127.0.0.1"
			},
			issue:       "The trusted proxy a127.0.0.1 is invalid",
			expectIssue: true,
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			testConfig = validConfig
			run.verifyFunc = testConfig.HTTP.Verify
			run.verifyIssuesAndError(t)
		})
	}
}

func TestConfig(t *testing.T) {
	var testConfig Config

	t.Run("expect no issues or error", func(t *testing.T) {
		testConfig := validConfig
		issues := testConfig.Verify()
		assert.Empty(t, issues, "expected no issues from verify function")
	})

	t.Run("expect multiple issues", func(t *testing.T) {
		testConfig := validConfig
		testConfig.LogLevel = ""
		testConfig.HTTP.ListenAddress = "aa"
		issues := testConfig.Verify()
		assert.Contains(t, issues, "An invalid log level was specified", "expected an issue but was not found")
		assert.Contains(t, issues, "HTTP listen address is not valid")
	})

	runs := []Run{
		// Log level
		{
			name:        "expect no log level error",
			beforeWork:  func() {},
			issue:       "An invalid log level was specified",
			expectIssue: false,
		},
		{
			name: "expect log level error when invalid log level is given",
			beforeWork: func() {
				testConfig.LogLevel = "a"
			},
			issue:       "An invalid log level was specified",
			expectIssue: true,
		},
		{
			name: "expect log level error when no LogLevel is not given",
			beforeWork: func() {
				testConfig.LogLevel = ""
			},
			issue:       "An invalid log level was specified",
			expectIssue: true,
		},
		// HTTP issues retrieved sanity check
		{
			name: "expect http issue to exist",
			beforeWork: func() {
				testConfig.HTTP.ListenAddress = ""
			},
			issue:       "HTTP listen address is not valid",
			expectIssue: true,
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			testConfig = validConfig
			run.verifyFunc = testConfig.Verify
			run.verifyIssuesAndError(t)
		})
	}
}
