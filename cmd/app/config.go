package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"strings"
)

type credential string

type config struct {
	BindingAddr string       `required:"false" envconfig:"BINDING_ADDRESS"`
	Debug       bool         `required:"false" envconfig:"DEBUG"`
	Credentials []credential `required:"false" envconfig:"CREDENTIALS"`

	RequestHeaderPrefix         string `required:"false" envconfig:"REQUEST_HEADER_PREFIX"`
	ForwardRequestHeaderPrefix  string `required:"false" envconfig:"FORWARD_REQUEST_HEADER_PREFIX"`
	ForwardResponseHeaderPrefix string `required:"false" envconfig:"FORWARD_RESPONSE_HEADER_PREFIX"`

	CorsAllowOrigin  []string `required:"false" envconfig:"CORS_ALLOW_ORIGIN"`
	CorsAllowMethods []string `required:"false" envconfig:"CORS_ALLOW_METHODS"`
	CorsAllowHeaders []string `required:"false" envconfig:"CORS_ALLOW_HEADERS"`
	CorsAllowMaxAge  int      `required:"false" envconfig:"CORS_ALLOW_MAX_AGE"`
}

func readConfig() (*config, error) {
	c := &config{
		BindingAddr:         ":8080",
		Debug:               false,
		RequestHeaderPrefix: "R-",
	}
	err := envconfig.Process("", c)

	if c.ForwardRequestHeaderPrefix == "" {
		c.ForwardRequestHeaderPrefix = c.RequestHeaderPrefix + "Forward-Request-Header-"
	}
	if c.ForwardResponseHeaderPrefix == "" {
		c.ForwardResponseHeaderPrefix = c.RequestHeaderPrefix + "Forward-Response-Header-"
	}

	if len(c.Credentials) == 0 {
		logrus.Warn("There are no user credentials configured!")
	}

	for _, c := range c.Credentials {
		if !strings.Contains(string(c), ":") {
			return nil, fmt.Errorf("invalid credential: %s", c)
		}
	}

	logrus.WithField("config", fmt.Sprintf("%#v", c)).Info("Active config")

	return c, err
}

func (c credential) UsernameAndPassword() (string, string) {
	split := strings.SplitN(string(c), ":", 2)

	return split[0], split[1]
}
