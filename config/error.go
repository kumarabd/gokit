// Copyright 2021 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/kumarabd/gokit/errors"
)

var (
	// ErrEmptyConfig is returned when the config has not been initialized.
	ErrEmptyConfig = errors.New("", errors.NoneSeverity, "Config not initialized")
)

// ErrViper returns a MeshKit error wrapping err in case of an (initialization) error in the Viper provider.
func ErrViper(err error) error {
	return errors.New("", errors.NoneSeverity, "Viper initialization failed with error", err.Error())
}

// ErrViper returns a MeshKit error wrapping err in case of an (initialization) error in the in-memory provider.
func ErrInMem(err error) error {
	return errors.New("", errors.NoneSeverity, "In Memory initialization failed with error", err.Error())
}
