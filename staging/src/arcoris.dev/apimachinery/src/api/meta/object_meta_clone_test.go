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

func TestObjectMetaClone(t *testing.T) {
	deletion := stamp.Deletion{DeletedAt: stamp.NewTimestamp(testTime())}
	meta := validObjectMeta()
	meta.Deletion = &deletion

	cloned := meta.Clone()
	cloned.Labels["role"] = "manager"
	cloned.Annotations["note"] = "changed"
	cloned.OwnerReferences[0].Controller = false
	cloned.Finalizers[0] = "changed"
	cloned.Deletion.GracePeriodSeconds = 99

	if meta.Labels["role"] != "worker" {
		t.Fatal("labels were not detached")
	}
	if meta.Annotations["note"] != "human readable" {
		t.Fatal("annotations were not detached")
	}
	if !meta.OwnerReferences[0].Controller {
		t.Fatal("owner references were not detached")
	}
	if meta.Finalizers[0] != "cleanup" {
		t.Fatal("finalizers were not detached")
	}
	if meta.Deletion.GracePeriodSeconds != 0 {
		t.Fatal("deletion pointer was not detached")
	}

	meta.Labels["role"] = "lead"
	meta.Annotations["note"] = "original changed"
	meta.OwnerReferences[0].Object.Name = "changed"
	meta.Finalizers[0] = "original"
	meta.Deletion.GracePeriodSeconds = 1

	if cloned.Labels["role"] != "manager" {
		t.Fatal("original label mutation changed clone")
	}
	if cloned.Annotations["note"] != "changed" {
		t.Fatal("original annotation mutation changed clone")
	}
	if cloned.OwnerReferences[0].Object.Name == "changed" {
		t.Fatal("original owner reference mutation changed clone")
	}
	if cloned.Finalizers[0] != "changed" {
		t.Fatal("original finalizer mutation changed clone")
	}
	if cloned.Deletion.GracePeriodSeconds != 99 {
		t.Fatal("original deletion mutation changed clone")
	}
}

func TestObjectMetaClonePreservesEmptyCollections(t *testing.T) {
	meta := ObjectMeta{
		Labels:          labels.Set{},
		Annotations:     annotations.Set{},
		OwnerReferences: owner.List{},
		Finalizers:      finalizer.Set{},
	}

	cloned := meta.Clone()
	if cloned.Labels == nil {
		t.Fatal("empty labels clone is nil")
	}
	if cloned.Annotations == nil {
		t.Fatal("empty annotations clone is nil")
	}
	if cloned.OwnerReferences == nil {
		t.Fatal("empty owner references clone is nil")
	}
	if cloned.Finalizers == nil {
		t.Fatal("empty finalizers clone is nil")
	}
	if cloned.Deletion != nil {
		t.Fatal("nil deletion clone is non-nil")
	}
}
