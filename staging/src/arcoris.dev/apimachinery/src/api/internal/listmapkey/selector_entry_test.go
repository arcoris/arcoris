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

package listmapkey

import (
	"math"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestExtractSelectorEntryMissingKey(t *testing.T) {
	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("status", value.StringValue("True")),
		objectElement(types.Field("type").String().Required()),
		[]types.FieldName{"type"},
		Options{},
	)

	requireErrorKind(t, err, FailureMissingKey)
	requireEqual(t, IsPayloadFailure(err), true)
}

func TestExtractSelectorEntryNullKey(t *testing.T) {
	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.NullValue()),
		objectElement(types.Field("type").String().Required()),
		[]types.FieldName{"type"},
		Options{},
	)

	requireErrorKind(t, err, FailureNullKey)
	requireEqual(t, IsPayloadFailure(err), true)
}

func TestExtractSelectorEntryWrongKeyKind(t *testing.T) {
	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.Int64Value(1)),
		objectElement(types.Field("type").String().Required()),
		[]types.FieldName{"type"},
		Options{},
	)

	requireErrorKind(t, err, FailureKeyKindMismatch)
	requireEqual(t, IsPayloadFailure(err), true)
}

func TestExtractSelectorIntegerRangeMismatch(t *testing.T) {
	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("port", value.Uint64Value(math.MaxUint64)),
		objectElement(types.Field("port").Int64().Required()),
		[]types.FieldName{"port"},
		Options{},
	)

	requireErrorKind(t, err, FailureKeyIntegerRange)
	requireEqual(t, IsPayloadFailure(err), true)
}

func TestExtractSelectorEntryRejectsNonObjectItem(t *testing.T) {
	_, err := ExtractSelector(
		conditionPath(0),
		value.StringValue("not-object"),
		objectElement(types.Field("type").String().Required()),
		[]types.FieldName{"type"},
		Options{},
	)

	requireErrorKind(t, err, FailureItemKindMismatch)
	requireEqual(t, IsPayloadFailure(err), true)
}
