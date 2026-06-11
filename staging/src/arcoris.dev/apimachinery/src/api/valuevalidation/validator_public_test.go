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

func TestValidatorValidate(t *testing.T) {
	validator := valuevalidation.New(valuevalidation.Options{})

	err := validator.Validate(value.StringValue("x"), types.Int64().Descriptor())

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$",
	)
}

func TestValidatorValidateAt(t *testing.T) {
	validator := valuevalidation.New(valuevalidation.Options{})
	path := rootField("desired", "replicas")

	err := validator.ValidateAt(path, value.StringValue("x"), types.Int64().Descriptor())

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$.desired.replicas",
	)
}

func TestValidatorCreatesFreshRunPerCall(t *testing.T) {
	validator := valuevalidation.New(valuevalidation.Options{MaxErrors: 1})
	shape := types.Object(
		types.Field("first").String().Required(),
		types.Field("second").String().Required(),
	).Descriptor()

	first := validator.Validate(mustObject(t), shape)
	second := validator.Validate(mustObject(t), shape)

	requireErrorCount(t, first, 1)
	requireErrorCount(t, second, 1)
}
