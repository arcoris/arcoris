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

func TestValidateMapAcceptsDynamicEntries(t *testing.T) {
	shape := types.MapOf(types.String().MinLen(1)).Type()
	payload := mustObject(
		t,
		value.ObjectMember("app", value.StringValue("api")),
		value.ObjectMember("tier", value.StringValue("backend")),
	)

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			shape,
			valuevalidation.Options{},
		),
	)
}

func TestValidateMapRejectsInvalidEntryValue(t *testing.T) {
	shape := types.MapOf(types.String().MinLen(1)).Type()
	payload := mustObject(t, value.ObjectMember("app", value.StringValue("")))

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		`$["app"]`,
	)
}

func TestValidateMapUsesKeyPath(t *testing.T) {
	shape := types.MapOf(types.Int32()).Type()
	payload := mustObject(t, value.ObjectMember("app.kubernetes.io/name", value.StringValue("api")))

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		`$["app.kubernetes.io/name"]`,
	)
}
