package config

import (
	"github.com/kumarabd/gokit/errors"
)

var (
	// ErrEmptyConfig is returned when the config has not been initialized.
	ErrEmptyConfig = errors.New("", errors.NoneSeverity, "Config not initialized")
)

// ErrViper returns a error wrapping err in case of an (initialization) error in the Viper provider.
func ErrViper(err error) error {
	return errors.New("", errors.NoneSeverity, "Viper initialization failed with error", err.Error())
}

// ErrInMem returns a error wrapping err in case of an (initialization) error in the in-memory provider.
func ErrInMem(err error) error {
	return errors.New("", errors.NoneSeverity, "In Memory initialization failed with error", err.Error())
}
