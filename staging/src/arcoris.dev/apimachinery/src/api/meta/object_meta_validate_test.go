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

package meta

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/finalizer"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/meta/owner"
	"arcoris.dev/apimachinery/api/meta/stamp"
)

func TestObjectMetaValidate(t *testing.T) {
	requireNoError(t, (ObjectMeta{}).Validate())
	requireNoError(t, validObjectMeta().Validate())
}

func TestObjectMetaValidateRejectsInvalidNestedMetadata(t *testing.T) {
	tests := []struct {
		name string
		meta ObjectMeta
	}{
		{name: "name", meta: ObjectMeta{Name: "Worker"}},
		{name: "generateName", meta: ObjectMeta{GenerateName: "-worker"}},
		{name: "namespace", meta: ObjectMeta{Namespace: "System"}},
		{name: "uid", meta: ObjectMeta{UID: "uid 1"}},
		{name: "resourceVersion", meta: ObjectMeta{ResourceVersion: "rv 1"}},
		{name: "labels", meta: ObjectMeta{Labels: labels.Set{"Role": "worker"}}},
		{name: "annotations", meta: ObjectMeta{Annotations: annotations.Set{"note": "bad\nnote"}}},
		{name: "owners", meta: ObjectMeta{OwnerReferences: owner.List{{}}}},
		{name: "finalizers", meta: ObjectMeta{Finalizers: finalizer.Set{"Cleanup"}}},
		{name: "deletion", meta: ObjectMeta{Deletion: &stamp.Deletion{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.meta.Validate(), ErrInvalidObjectMeta)
		})
	}
}
