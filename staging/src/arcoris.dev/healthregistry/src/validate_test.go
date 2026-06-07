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

package healthregistry

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestPrepareChecksJoinsBatchValidationFailures(t *testing.T) {
	t.Parallel()

	var typed *typedNilChecker
	prepared, err := prepareChecks(
		health.TargetReady,
		[]health.Checker{
			nil,
			typed,
			namedChecker{name: "bad-name"},
			checker("dup"),
			checker("dup"),
		},
	)

	if prepared != nil {
		t.Fatalf("prepared = %+v, want nil on validation failure", prepared)
	}
	if !errors.Is(err, health.ErrNilChecker) ||
		!errors.Is(err, health.ErrInvalidCheckName) ||
		!errors.Is(err, ErrDuplicateCheck) {
		t.Fatalf("prepareChecks() = %v, want joined nil/name/duplicate errors", err)
	}
}
