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

package healthgrpc

import (
	"errors"
	"testing"
	"time"
)

func TestInvalidWatchIntervalError(t *testing.T) {
	t.Parallel()

	err := error(InvalidWatchIntervalError{Interval: -time.Second})
	if !errors.Is(err, ErrInvalidWatchInterval) {
		t.Fatalf("errors.Is(%v, ErrInvalidWatchInterval) = false, want true", err)
	}

	var intervalErr InvalidWatchIntervalError
	if !errors.As(err, &intervalErr) {
		t.Fatalf("errors.As(%T, InvalidWatchIntervalError) = false, want true", err)
	}
	if intervalErr.Interval != -time.Second {
		t.Fatalf("Interval = %s, want -1s", intervalErr.Interval)
	}
	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}

func TestInvalidMaxListServicesError(t *testing.T) {
	t.Parallel()

	err := error(InvalidMaxListServicesError{Max: -1})
	if !errors.Is(err, ErrInvalidMaxListServices) {
		t.Fatalf("errors.Is(%v, ErrInvalidMaxListServices) = false, want true", err)
	}

	var maxErr InvalidMaxListServicesError
	if !errors.As(err, &maxErr) {
		t.Fatalf("errors.As(%T, InvalidMaxListServicesError) = false, want true", err)
	}
	if maxErr.Max != -1 {
		t.Fatalf("Max = %d, want -1", maxErr.Max)
	}
	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}
