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
	shape := types.MapOf(types.String().MinBytes(1)).Descriptor()
	payload := mustObject(
		t,
		value.MustRecordMember("app", value.StringValue("api")),
		value.MustRecordMember("tier", value.StringValue("backend")),
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
	shape := types.MapOf(types.String().MinBytes(1)).Descriptor()
	payload := mustObject(t, value.MustRecordMember("app", value.StringValue("")))

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

func TestValidateMapRejectsInvalidEntryKey(t *testing.T) {
	shape := types.MapOf(types.String()).
		Keys(types.String().Pattern(`^[a-z]+$`)).
		Descriptor()
	payload := mustObject(t, value.MustRecordMember("INVALID", value.StringValue("ok")))

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrPatternMismatch,
		valuevalidation.ErrorReasonPatternMismatch,
		`$["INVALID"]`,
	)
}

func TestValidateMapRejectsInvalidEntryKeyLength(t *testing.T) {
	shape := types.MapOf(types.String()).
		Keys(types.String().MinBytes(2)).
		Descriptor()
	payload := mustObject(t, value.MustRecordMember("a", value.StringValue("ok")))

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
		`$["a"]`,
	)
}

func TestValidateMapValidatesEntryValueWhenKeyIsInvalid(t *testing.T) {
	shape := types.MapOf(types.Int32()).
		Keys(types.String().MinBytes(2)).
		Descriptor()
	payload := mustObject(t, value.MustRecordMember("a", value.StringValue("not-int")))

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireErrorCount(t, err, 2)
	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		`$["a"]`,
	)
	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		`$["a"]`,
	)
}

func TestValidateMapUsesKeyPath(t *testing.T) {
	shape := types.MapOf(types.Int32()).Descriptor()
	payload := mustObject(t, value.MustRecordMember("app.kubernetes.io/name", value.StringValue("api")))

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
