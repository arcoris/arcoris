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

package valuevalidation_test

import (
	"math"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateUnsignedIntegerConstraints(t *testing.T) {
	tests := []struct {
		name     string
		payload  value.Value
		shape    types.Descriptor
		sentinel error
		reason   valuevalidation.ErrorReason
	}{
		{
			name:     "negative value",
			payload:  value.Int64Value(-1),
			shape:    types.Uint64().Descriptor(),
			sentinel: valuevalidation.ErrValueOutOfRange,
			reason:   valuevalidation.ErrorReasonBelowMinimum,
		},
		{
			name:     "uint8 overflow",
			payload:  value.Uint64Value(math.MaxUint8 + 1),
			shape:    types.Uint8().Descriptor(),
			sentinel: valuevalidation.ErrValueOutOfRange,
			reason:   valuevalidation.ErrorReasonAboveMaximum,
		},
		{
			name:     "uint32 enum mismatch",
			payload:  value.Uint64Value(3),
			shape:    types.Uint32().Enum(1, 2).Descriptor(),
			sentinel: valuevalidation.ErrEnumMismatch,
			reason:   valuevalidation.ErrorReasonEnumMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valuevalidation.Validate(
				tt.payload,
				tt.shape,
				valuevalidation.Options{},
			)

			requireError(t, err, tt.sentinel, tt.reason, "$")
		})
	}
}

func TestValidateUint64AcceptsMaxUint64(t *testing.T) {
	requireNoError(
		t,
		valuevalidation.Validate(
			value.Uint64Value(math.MaxUint64),
			types.Uint64().Descriptor(),
			valuevalidation.Options{},
		),
	)
}
