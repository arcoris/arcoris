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
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestSelectorStringKey(t *testing.T) {
	gotSelector, err := Selector(
		conditionPath(0),
		objectWith("type", value.StringValue("Ready")),
		objectElement(types.Field("type").String().Required()),
		[]types.FieldName{"type"},
		Options{},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"type":"Ready"}`)
}

func TestSelectorBoolKey(t *testing.T) {
	gotSelector, err := Selector(
		conditionPath(0),
		objectWith("enabled", value.BoolValue(true)),
		objectElement(types.Field("enabled").Bool().Required()),
		[]types.FieldName{"enabled"},
		Options{},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"enabled":true}`)
}

func TestSelectorSignedIntegerKey(t *testing.T) {
	gotSelector, err := Selector(
		conditionPath(0),
		objectWith("port", value.Int64Value(-1)),
		objectElement(types.Field("port").Int64().Required()),
		[]types.FieldName{"port"},
		Options{},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"port":-1}`)
}

func TestSelectorUnsignedIntegerKey(t *testing.T) {
	gotSelector, err := Selector(
		conditionPath(0),
		objectWith("port", value.Uint64Value(443)),
		objectElement(types.Field("port").Uint64().Required()),
		[]types.FieldName{"port"},
		Options{},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"port":443}`)
}

func TestSelectorMultiKeyCanonicalSelector(t *testing.T) {
	gotSelector, err := Selector(
		conditionPath(0),
		objectWithMembers(
			value.ObjectMember("port", value.Uint64Value(443)),
			value.ObjectMember("host", value.StringValue("api.example.com")),
		),
		objectElement(
			types.Field("port").Uint64().Required(),
			types.Field("host").String().Required(),
		),
		[]types.FieldName{"port", "host"},
		Options{},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"host":"api.example.com","port":443}`)
}
