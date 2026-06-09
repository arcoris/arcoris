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
	"testing"
	"time"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateTemporalAcceptsMatchingKinds(t *testing.T) {
	timeOfDay, err := value.NewTimeOfDay(12, 30, 0, 0)
	if err != nil {
		t.Fatalf("NewTimeOfDay() error = %v", err)
	}
	timeOfDayValue, err := value.TimeOfDayValue(timeOfDay)
	if err != nil {
		t.Fatalf("TimeOfDayValue() error = %v", err)
	}

	tests := []struct {
		name    string
		payload value.Value
		shape   types.Descriptor
	}{
		{
			name:    "timestamp",
			payload: value.TimestampValue(time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)),
			shape:   types.Timestamp().Descriptor(),
		},
		{
			name:    "date",
			payload: mustDate(t, 2026, 6, 1),
			shape:   types.Date().Descriptor(),
		},
		{
			name:    "time of day",
			payload: timeOfDayValue,
			shape:   types.Time().Descriptor(),
		},
		{
			name:    "duration",
			payload: value.DurationValue(time.Second),
			shape:   types.Duration().Descriptor(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireNoError(
				t,
				valuevalidation.Validate(
					tt.payload,
					tt.shape,
					valuevalidation.Options{},
				),
			)
		})
	}
}

func TestValidateTemporalRejectsMismatchedKind(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("2026-06-01"),
		types.Date().Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$",
	)
}
