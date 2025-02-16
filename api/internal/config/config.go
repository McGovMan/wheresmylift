package config

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
)

type HTTP struct {
	ListenAddress string `mapstructure:"listen_address" yaml:"listen_address"`
	TrustedProxy  string `mapstructure:"trusted_proxy" yaml:"trusted_proxy"`
}

// Config describes the configuration for Server
type Config struct {
	LogLevel string `mapstructure:"log_level" yaml:"log_level"`
	HTTP     HTTP   `mapstructure:"http" yaml:"http"`
}

func (c *Config) GetZeroLogLevel() zerolog.Level {
	switch c.LogLevel {
	case "trace":
		return zerolog.TraceLevel
	case "disabled":
		return zerolog.Disabled
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.NoLevel
	}
}

func (h *HTTP) Verify() []string {
	issues := []string{}
	_, _, err := net.SplitHostPort(h.ListenAddress)
	if err != nil {
		issues = append(issues, "HTTP listen address is not valid")
	}

	if net.ParseIP(h.TrustedProxy) == nil {
		issues = append(issues, fmt.Sprintf("The trusted proxy %s is invalid", h.TrustedProxy))
	}

	return issues
}

func (c *Config) Verify() []string {
	issues := []string{}

	if c.GetZeroLogLevel() == zerolog.NoLevel {
		issues = append(issues, "An invalid log level was specified")
	}

	httpIssues := c.HTTP.Verify()
	issues = append(issues, httpIssues...)

	return issues
}
