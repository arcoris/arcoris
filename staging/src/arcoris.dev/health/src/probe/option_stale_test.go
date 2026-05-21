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

package probe

import (
	"errors"
	"testing"
	"time"
)

func TestWithStaleAfter(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithStaleAfter(0)(&cfg)
	if err != nil {
		t.Fatalf("WithStaleAfter(0) = %v, want nil", err)
	}
	if cfg.staleAfter != 0 {
		t.Fatalf("staleAfter = %s, want 0", cfg.staleAfter)
	}
}

func TestWithStaleAfterRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithStaleAfter(-time.Nanosecond)(&cfg)

	if !errors.Is(err, ErrInvalidStaleAfter) {
		t.Fatalf("WithStaleAfter(-1ns) = %v, want ErrInvalidStaleAfter", err)
	}
}
