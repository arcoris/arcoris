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
	"strings"
	"testing"
)

func TestInvalidConfigError(t *testing.T) {
	t.Parallel()

	err := invalidConfig("decode.limits.max_depth", "must be greater than zero")
	requireConfigErrorIs(t, err, ErrInvalidConfig)
	requireErrorTextContains(t, err, "decode.limits.max_depth")
	requireErrorTextContains(t, err, "must be greater than zero")
}

func TestUnsupportedConfigError(t *testing.T) {
	t.Parallel()

	err := unsupportedConfig("encode.values.bytes", "mode is not implemented")
	requireConfigErrorIs(t, err, ErrUnsupportedConfig)
	requireErrorTextContains(t, err, "encode.values.bytes")
	requireErrorTextContains(t, err, "mode is not implemented")
}

func TestNilConfigError(t *testing.T) {
	t.Parallel()

	var err *configError

	if got := err.Error(); got != "<nil>" {
		t.Fatalf("nil Error() = %q; want %q", got, "<nil>")
	}
	if got := err.Unwrap(); got != nil {
		t.Fatalf("nil Unwrap() = %v; want nil", got)
	}
}

func requireConfigErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(..., %v)", err, target)
	}
}

func requireErrorTextContains(t *testing.T, err error, want string) {
	t.Helper()

	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %v; want text containing %q", err, want)
	}
}
