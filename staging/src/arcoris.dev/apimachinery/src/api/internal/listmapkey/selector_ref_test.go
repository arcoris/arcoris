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

func TestExtractSelectorRefToObjectElement(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		if name == "example.Condition" {
			return types.Define(
				"example.Condition",
				types.Object(types.Field("type").String().Required()),
			), true
		}

		return types.TypeDefinition{}, false
	})

	gotSelector, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.StringValue("Ready")),
		types.Ref("example.Condition").Type(),
		[]types.FieldName{"type"},
		Options{Resolver: typeResolver},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"type":"Ready"}`)
}

func TestExtractSelectorRefKeyType(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		if name == "example.ConditionType" {
			return types.Define("example.ConditionType", types.String()), true
		}

		return types.TypeDefinition{}, false
	})

	gotSelector, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.StringValue("Ready")),
		objectElement(types.Field("type").Ref("example.ConditionType").Required()),
		[]types.FieldName{"type"},
		Options{Resolver: typeResolver},
	)

	requireNoError(t, err)
	requireEqual(t, gotSelector.String(), `{"type":"Ready"}`)
}

func TestExtractSelectorMissingResolver(t *testing.T) {
	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.StringValue("Ready")),
		types.Ref("example.Condition").Type(),
		[]types.FieldName{"type"},
		Options{},
	)

	requireErrorKind(t, err, FailureUnresolvedRef)
	requireEqual(t, IsDescriptorFailure(err), true)
}

func TestExtractSelectorUnresolvedRef(t *testing.T) {
	typeResolver := resolverFunc(func(types.TypeName) (types.TypeDefinition, bool) {
		return types.TypeDefinition{}, false
	})

	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.StringValue("Ready")),
		types.Ref("example.Condition").Type(),
		[]types.FieldName{"type"},
		Options{Resolver: typeResolver},
	)

	requireErrorKind(t, err, FailureUnresolvedRef)
	requireEqual(t, IsDescriptorFailure(err), true)
}

func TestExtractSelectorReferenceCycle(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		if name == "example.Condition" {
			return types.Define("example.Condition", types.Ref("example.Condition")), true
		}

		return types.TypeDefinition{}, false
	})

	_, err := ExtractSelector(
		conditionPath(0),
		objectWith("type", value.StringValue("Ready")),
		types.Ref("example.Condition").Type(),
		[]types.FieldName{"type"},
		Options{Resolver: typeResolver},
	)

	requireErrorKind(t, err, FailureReferenceCycle)
	requireEqual(t, IsDescriptorFailure(err), true)
}
