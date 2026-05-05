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

package healthprobe

import (
	"errors"
	"testing"
	"time"
)

func TestInvalidStaleAfterErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := InvalidStaleAfterError{StaleAfter: -time.Second}

	if !errors.Is(err, ErrInvalidStaleAfter) {
		t.Fatalf("errors.Is(%v, ErrInvalidStaleAfter) = false, want true", err)
	}
}

func TestInvalidStaleAfterErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidStaleAfterError{StaleAfter: -time.Second})

	var staleErr InvalidStaleAfterError
	if !errors.As(err, &staleErr) {
		t.Fatalf("errors.As(%T, InvalidStaleAfterError) = false, want true", err)
	}
	if staleErr.StaleAfter != -time.Second {
		t.Fatalf("StaleAfter = %s, want %s", staleErr.StaleAfter, -time.Second)
	}
}

func TestInvalidStaleAfterErrorMessage(t *testing.T) {
	t.Parallel()

	err := InvalidStaleAfterError{StaleAfter: -time.Second}

	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}
