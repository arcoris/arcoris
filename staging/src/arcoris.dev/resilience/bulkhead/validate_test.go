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

package bulkhead

import (
	"errors"
	"testing"
)

func TestValidateConfigRejectsZeroLimit(t *testing.T) {
	t.Parallel()

	cfg := newConfig(0)
	if err := validateConfig(cfg); !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("validateConfig error = %v, want %v", err, ErrInvalidLimit)
	}
}

func TestValidateConfigRejectsNilClock(t *testing.T) {
	t.Parallel()

	cfg := newConfig(1)
	cfg.clock = nil
	if err := validateConfig(cfg); !errors.Is(err, ErrNilClock) {
		t.Fatalf("validateConfig error = %v, want %v", err, ErrNilClock)
	}
}

func TestValidateConfigAcceptsValidConfig(t *testing.T) {
	t.Parallel()

	if err := validateConfig(newConfig(1)); err != nil {
		t.Fatalf("validateConfig returned error: %v", err)
	}
}
