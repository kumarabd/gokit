package viper

import (
	"fmt"
	"sync"

	"github.com/kumarabd/gokit/config"
	"github.com/spf13/viper"
)

const (
	// ViperKey corresponds to viper configuration
	ViperKey = "viper"
	// FilePath the path to store the config file
	FilePath = "filepath"
	// FileType the type of the config file
	FileType = "filetype"
	// FileName the name of the config file
	FileName = "filename"
)

// Options contains config options for various aspects of viper instance.
type Options struct {
	FilePath string
	FileType string
	FileName string
}

// Viper implements the config interface Handler for a Viper configuration registry.
type Viper struct {
	instance *viper.Viper
	mutex    sync.Mutex
}

// NewViper returns a new instance of a Viper configuration provider using the provided Options opts.
func New(opts Options) (config.Handler, error) {
	v := viper.New()
	v.AddConfigPath(opts.FilePath)
	v.SetConfigType(opts.FileType)
	v.SetConfigName(opts.FileName)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
			// Hack until viper issue #433 is fixed
			er := v.WriteConfigAs(fmt.Sprintf("%s/%s.%s", opts.FilePath, opts.FileName, opts.FileType))
			if er != nil {
				return nil, config.ErrViper(err)
			}
			_ = v.WriteConfig()
		} else {
			// Config file was found but another error was produced
			return nil, config.ErrViper(err)
		}
	}

	return &Viper{
		instance: v,
	}, nil
}

// SetKey implements SetKey functionality of the interface
func (v *Viper) SetKey(key string, value string) {
	v.mutex.Lock()
	v.instance.Set(key, value)
	_ = v.instance.WriteConfig()
	v.mutex.Unlock()
}

// GetKey implements GetKey functionality of the interface
func (v *Viper) GetKey(key string) string {
	v.mutex.Lock()
	_ = v.instance.ReadInConfig()
	defer v.mutex.Unlock()
	return v.instance.Get(key).(string)
}

// GetObject implements GetObject functionality of the interface
func (v *Viper) GetObject(key string, result interface{}) error {
	v.mutex.Lock()
	_ = v.instance.ReadInConfig()
	err := v.instance.UnmarshalKey(key, &result)
	defer v.mutex.Unlock()
	if err != nil {
		return err
	}
	return err
}

// SetObject implements SetObject functionality of the interface
func (v *Viper) SetObject(key string, value interface{}) error {
	v.mutex.Lock()
	v.instance.Set(key, value)
	err := v.instance.WriteConfig()
	defer v.mutex.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// GetAll implements GetAll functionality of the interface
func (v *Viper) GetAll(result interface{}) error {
	v.mutex.Lock()
	_ = v.instance.ReadInConfig()
	err := v.instance.Unmarshal(&result)
	defer v.mutex.Unlock()
	if err != nil {
		return err
	}
	return err
}

// Is implements Is functionality of the interface
func (v *Viper) Is(key string) bool {
	return v.instance.IsSet(key)
}
