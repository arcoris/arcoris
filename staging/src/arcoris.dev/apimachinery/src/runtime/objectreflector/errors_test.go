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

package objectreflector

import (
	"errors"
	"testing"
)

func requireNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func TestErrorBuildersPreserveSentinelsAndCauses(t *testing.T) {
	cause := errors.New("cause")
	tests := []struct {
		name   string
		err    error
		target error
	}{
		{name: "invalid option", err: invalidOptionError(cause), target: ErrInvalidOption},
		{name: "invalid event", err: invalidEventError(cause), target: ErrInvalidEvent},
		{name: "outside collection", err: changeOutsideCollectionError(cause), target: ErrChangeOutsideCollection},
		{name: "non monotonic revision", err: nonMonotonicRevisionError(cause), target: ErrNonMonotonicRevision},
		{name: "source contract", err: sourceContractError("source failed: %w", cause), target: ErrInvalidEvent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, tt.target)
			requireErrorIs(t, tt.err, cause)
		})
	}
}
