// Copyright 2026 The ARCORIS Authors.
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

package jsonconfig

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfig classifies structurally invalid JSON codec config.
	ErrInvalidConfig = errors.New("invalid JSON codec config")

	// ErrUnsupportedConfig classifies known-but-unimplemented JSON codec config.
	ErrUnsupportedConfig = errors.New("unsupported JSON codec config")
)

// configError records one invalid or unsupported config field.
type configError struct {
	// kind is the sentinel exposed through errors.Is.
	kind error

	// path identifies the public config field using dotted lower-case names.
	path string

	// detail explains the invalid or unsupported value.
	detail string
}

// Error returns a stable config diagnostic.
func (e *configError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%s: %s: %s", e.kind, e.path, e.detail)
}

// Unwrap exposes the broad config sentinel.
func (e *configError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.kind
}

// invalidConfig reports an invalid config field.
func invalidConfig(path string, format string, args ...any) error {
	return &configError{
		kind:   ErrInvalidConfig,
		path:   path,
		detail: fmt.Sprintf(format, args...),
	}
}

// unsupportedConfig reports a known config field that codecjson does not implement.
func unsupportedConfig(path string, format string, args ...any) error {
	return &configError{
		kind:   ErrUnsupportedConfig,
		path:   path,
		detail: fmt.Sprintf(format, args...),
	}
}
