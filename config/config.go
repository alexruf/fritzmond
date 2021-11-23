package config

import (
	"errors"
)

type Config struct {
	Url           string
	Username      string
	Password      string
	SkipTlsVerify bool
	Interval      uint
	DbPath        string
	Port          uint
	DisableWebUi  bool
}

// Validate returns and error in case the values in Config violate some basic validation rules.
func (c Config) Validate() error {
	if len(c.Url) == 0 {
		return errors.New("no value set for URL")
	}
	if len(c.Username) == 0 {
		return errors.New("no value set for username")
	}
	if len(c.Password) == 0 {
		return errors.New("no value set for password")
	}
	return nil
}
