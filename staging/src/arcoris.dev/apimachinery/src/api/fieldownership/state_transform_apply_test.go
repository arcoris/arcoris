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

package fieldownership

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestTransformOwnerFields(t *testing.T) {
	state, err := baseState().transformOwnerFields(
		owner("user-cli"),
		set(namePath()),
		"bad path",
		unionTransform,
	)

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.metadata.name", "$.spec.image", "$.spec.replicas")
}

func TestTransformOtherOwnerFields(t *testing.T) {
	state, err := baseState().transformOtherOwnerFields(
		owner("user-cli"),
		set(replicasPath()),
		"bad path",
		removeExactTransform,
	)

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image", "$.spec.replicas")
	requireSet(t, state.FieldsFor(owner("autoscaler")))
}

func TestReplaceOwnerFields(t *testing.T) {
	state, err := baseState().replaceOwnerFields(owner("user-cli"), set(namePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.metadata.name")
}

func TestReplaceOwnerFieldsEmptyRemovesOwner(t *testing.T) {
	state, err := baseState().replaceOwnerFields(owner("user-cli"), fieldpath.EmptySet())

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")))
}
