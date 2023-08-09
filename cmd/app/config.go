package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"strings"
)

type credential string

type config struct {
	BindingAddr           string       `required:"false" envconfig:"BINDING_ADDRESS"`
	Debug                 bool         `required:"false" envconfig:"DEBUG"`
	RequestHeaderPrefix   string       `required:"false" envconfig:"REQUEST_HEADER_PREFIX"`
	RequestCredentials    []credential `required:"true" envconfig:"REQUEST_CREDENTIALS"`
	ForwardRequestHeader  string       `required:"false" envconfig:"FORWARD_REQUEST_HEADER"`
	ForwardResponseHeader string       `required:"false" envconfig:"FORWARD_RESPONSE_HEADER"`
}

func readConfig() (*config, error) {
	c := &config{
		BindingAddr:         ":8080",
		Debug:               false,
		RequestHeaderPrefix: "R-",
	}
	err := envconfig.Process("", c)

	if c.ForwardRequestHeader == "" {
		c.ForwardRequestHeader = c.RequestHeaderPrefix + "Forward-Request-Header"
	}
	if c.ForwardResponseHeader == "" {
		c.ForwardResponseHeader = c.RequestHeaderPrefix + "Forward-Response-Header"
	}

	for _, c := range c.RequestCredentials {
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
