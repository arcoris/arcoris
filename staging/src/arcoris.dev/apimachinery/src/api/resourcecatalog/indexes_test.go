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

package resourcecatalog

import (
	"testing"

	"arcoris.dev/apimachinery/api/identity"
)

func TestIndexKeyHelpers(t *testing.T) {
	def := validDefinition(
		"Worker",
		"workers",
		objectVersion("v1alpha1"),
		objectVersion("v1"),
	)

	requireEqual(t, groupResourceOf(def), identity.GroupResource{Group: testGroup, Resource: "workers"})
	requireEqual(t, groupKindOf(def), identity.GroupKind{Group: testGroup, Kind: "Worker"})

	requireSliceEqual(
		t,
		versionResourceKeys(def),
		[]identity.GroupVersionResource{
			{Group: testGroup, Version: "v1alpha1", Resource: "workers"},
			{Group: testGroup, Version: "v1", Resource: "workers"},
		},
	)
	requireSliceEqual(
		t,
		versionKindKeys(def),
		[]identity.GroupVersionKind{
			{Group: testGroup, Version: "v1alpha1", Kind: "Worker"},
			{Group: testGroup, Version: "v1", Kind: "Worker"},
		},
	)
}

func TestIndexPathHelpers(t *testing.T) {
	requireEqual(t, definitionPath(3), "definitions[3]")
	requireEqual(
		t,
		resourcePath(identity.GroupResource{
			Group:    testGroup,
			Resource: "workers",
		}),
		"definitions[control.arcoris.dev:workers]",
	)
	requireEqual(
		t,
		kindPath(identity.GroupKind{
			Group: testGroup,
			Kind:  "Worker",
		}),
		"definitions[control.arcoris.dev#Worker]",
	)
	requireEqual(
		t,
		versionResourcePath(identity.GroupVersionResource{
			Group:    testGroup,
			Version:  "v1",
			Resource: "workers",
		}),
		"definitions[control.arcoris.dev/v1:workers]",
	)
	requireEqual(
		t,
		versionKindPath(identity.GroupVersionKind{
			Group:   testGroup,
			Version: "v1",
			Kind:    "Worker",
		}),
		"definitions[control.arcoris.dev/v1#Worker]",
	)
}
