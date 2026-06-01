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

package valuevalidation

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestValidateIntegerLimitsReportsBounds(t *testing.T) {
	tests := []struct {
		name   string
		got    int64
		limits integerLimits[int64]
		reason ErrorReason
	}{
		{
			name: "below minimum",
			got:  1,
			limits: integerLimits[int64]{
				lower: integerBound[int64]{value: 2, set: true},
			},
			reason: ErrorReasonBelowMinimum,
		},
		{
			name: "above maximum",
			got:  3,
			limits: integerLimits[int64]{
				upper: integerBound[int64]{value: 2, set: true},
			},
			reason: ErrorReasonAboveMaximum,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run := newValidator(Options{})
			validateIntegerLimits(run, fieldpath.RootPath(), tt.got, tt.limits)

			err := run.result()
			requireInternalError(t, err, ErrValueOutOfRange, tt.reason, "$")
		})
	}
}

func TestValidateIntegerLimitsAcceptsValueWithinBounds(t *testing.T) {
	run := newValidator(Options{})

	validateIntegerLimits(
		run,
		fieldpath.RootPath(),
		int64(2),
		integerLimits[int64]{
			lower: integerBound[int64]{value: 1, set: true},
			upper: integerBound[int64]{value: 3, set: true},
		},
	)

	if err := run.result(); err != nil {
		t.Fatalf("result() = %v, want nil", err)
	}
}

func TestIntegerLimitsUnwrapSentinel(t *testing.T) {
	run := newValidator(Options{})
	validateIntegerLimits(
		run,
		fieldpath.RootPath(),
		int64(0),
		integerLimits[int64]{
			lower: integerBound[int64]{value: 1, set: true},
		},
	)

	if !errors.Is(run.result(), ErrValueOutOfRange) {
		t.Fatalf("errors.Is(ErrValueOutOfRange) = false")
	}
}
