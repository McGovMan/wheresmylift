package config

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/rs/zerolog"
)

type Timeouts struct {
	Startup    time.Duration `mapstructure:"startup" yaml:"startup"`
	Shutdown   time.Duration `mapstructure:"shutdown" yaml:"shutdown"`
	ReadHeader time.Duration `mapstructure:"read_header" yaml:"read_header"`
}

type CORS struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" yaml:"allowed_origins"`
}

type HTTP struct {
	ListenAddress  string `mapstructure:"listen_address" yaml:"listen_address"`
	CORS           CORS
	TrustedProxies []string `mapstructure:"trusted_proxies" yaml:"trusted_proxies"`
}

// Config describes the configuration for Server
type Config struct {
	LogLevel string   `mapstructure:"log_level" yaml:"log_level"`
	Timeouts Timeouts `mapstructure:"timeouts" yaml:"timeouts"`
	HTTP     HTTP     `mapstructure:"http" yaml:"http"`
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

func (t *Timeouts) Verify() []string {
	issues := []string{}
	timeoutsRegex := regexp.MustCompile("^[0-9]{1,2}s$")

	if !timeoutsRegex.MatchString(t.Startup.String()) || t.Startup.String() == "0s" {
		issues = append(issues, "Startup timeout should be represented in the form '{int}s', e.g. '30s'")
	}

	if !timeoutsRegex.MatchString(t.Shutdown.String()) || t.Shutdown.String() == "0s" {
		issues = append(issues, "Shutdown timeout should be represented in the form '{int}s', e.g. '30s'")
	}

	if !timeoutsRegex.MatchString(t.ReadHeader.String()) || t.ReadHeader.String() == "0s" {
		issues = append(issues, "ReadHeader timeout should be represented in the form '{int}s', e.g. '2s'")
	}

	return issues
}

func (h *HTTP) Verify() []string {
	issues := []string{}
	_, _, err := net.SplitHostPort(h.ListenAddress)
	if err != nil {
		issues = append(issues, "HTTP listen address is not valid")
	}

	if len(h.CORS.AllowedOrigins) == 0 {
		issues = append(issues, "No allowed origins specified")
	}

	httpDomainRegex := regexp.MustCompile("^(?:https?://)?([a-z0-9_-]+).+[a-z]{2,}$")
	for _, origin := range h.CORS.AllowedOrigins {
		if origin != "*" && !httpDomainRegex.MatchString(origin) {
			issues = append(issues, fmt.Sprintf("The allowed origin %s is invalid", origin))
		}
	}

	for _, proxy := range h.TrustedProxies {
		if net.ParseIP(proxy) == nil {
			issues = append(issues, fmt.Sprintf("The trusted proxy %s is invalid", proxy))
		}
	}

	return issues
}

func (c *Config) Verify() []string {
	issues := []string{}

	if c.GetZeroLogLevel() == zerolog.NoLevel {
		issues = append(issues, "An invalid log level was specified")
	}

	timeoutIssues := c.Timeouts.Verify()
	issues = append(issues, timeoutIssues...)

	httpIssues := c.HTTP.Verify()
	issues = append(issues, httpIssues...)

	return issues
}
