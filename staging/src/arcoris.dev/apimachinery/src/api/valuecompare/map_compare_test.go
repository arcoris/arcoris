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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"testing"
)

func TestCompareMapSameIsEmpty(t *testing.T) {
	descriptor := types.MapOf(types.String()).Descriptor()
	oldValue := valueObject("env", "prod")

	got, err := Compare(oldValue, oldValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareMapAddedKey(t *testing.T) {
	path := rootField("labels")

	got, err := CompareAt(path, valueObject(), valueObject("env", "prod"), types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Key("env")), nil, nil)
}

func TestCompareMapRemovedKey(t *testing.T) {
	path := rootField("labels")

	got, err := CompareAt(path, valueObject("env", "prod"), valueObject(), types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Key("env")), nil)
}

func TestCompareMapModifiedKeyValue(t *testing.T) {
	path := rootField("labels")

	got, err := CompareAt(path, valueObject("env", "prod"), valueObject("env", "stage"), types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Key("env")))
}

func TestCompareMapNestedObjectValue(t *testing.T) {
	path := rootField("ports")
	descriptor := types.MapOf(
		types.Object(types.Field("target").String().Optional()),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("http", valueObject("target", "8080")))
	newValue := value.MustRecordValue(value.MustRecordMember("http", valueObject("target", "8081")))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Key("http").Field("target")))
}

func TestCompareMapEmptyToNonEmpty(t *testing.T) {
	path := rootField("labels")

	got, err := CompareAt(path, valueObject(), valueObject("env", "prod"), types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Key("env")), nil, nil)
}

func TestCompareMapNonEmptyToEmpty(t *testing.T) {
	path := rootField("labels")

	got, err := CompareAt(path, valueObject("env", "prod"), valueObject(), types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Key("env")), nil)
}
func TestMapOperandReportsPresence(t *testing.T) {
	members := map[string]value.Value{"env": value.StringValue("prod")}

	got := mapOperand(members, "env")
	val, ok := got.ValueOK()
	text, _ := val.AsString()
	if !ok || text != "prod" {
		t.Fatalf("mapOperand(existing) = %#v", got)
	}

	got = mapOperand(members, "missing")
	if got.Present() {
		t.Fatalf("mapOperand(missing).Present() = true")
	}
}
