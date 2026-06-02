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

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateDecimalRejectsKindMismatch(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("12.30"),
		types.Decimal().Type(),
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

func TestValidateDecimalPrecisionAndScale(t *testing.T) {
	err := valuevalidation.Validate(
		mustDecimal(t, "123.45"),
		types.Decimal().Precision(4).Scale(1).Type(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrValueOutOfRange,
		valuevalidation.ErrorReasonAboveMaximum,
		"$",
	)
	requireErrorCount(t, err, 2)
}
