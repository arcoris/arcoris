/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfigtest

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidName reports that a test configuration name is empty.
	ErrInvalidName = errors.New("liveconfigtest: invalid config name")

	// ErrInvalidVersion reports that a test configuration version is negative.
	ErrInvalidVersion = errors.New("liveconfigtest: invalid config version")

	// ErrInvalidTimeout reports that a test configuration timeout is not positive.
	ErrInvalidTimeout = errors.New("liveconfigtest: invalid config timeout")

	// ErrInvalidLimit reports that a test configuration limit is negative.
	ErrInvalidLimit = errors.New("liveconfigtest: invalid config limit")
)

// ValidateConfig validates cfg according to liveconfigtest's fixture rules.
//
// The function is intentionally small and deterministic. It is intended for
// tests that need a stable validator for accepted and rejected live
// configuration updates.
func ValidateConfig(cfg Config) error {
	return cfg.validate()
}

// validate contains the fixture validation rules used by Config.Validate and
// ValidateConfig.
func (cfg Config) validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return ErrInvalidName
	}
	if cfg.Version < 0 {
		return ErrInvalidVersion
	}
	if cfg.Timeout <= 0 {
		return ErrInvalidTimeout
	}
	for _, limit := range cfg.Limits {
		if limit < 0 {
			return ErrInvalidLimit
		}
	}
	return nil
}
