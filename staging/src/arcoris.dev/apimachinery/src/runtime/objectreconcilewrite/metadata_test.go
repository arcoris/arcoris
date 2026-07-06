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

package objectreconcilewrite

import "testing"

func TestPatchMetadataBuildsRequestFromCurrent(t *testing.T) {
	key := testKey("task-1")
	owner := testOwner()
	current := newCurrent(t, key, 7, "desired")
	label := "api"
	annotation := "owned"

	req, err := current.PatchMetadata(
		map[string]*string{"app": &label},
		map[string]*string{"note": &annotation},
		owner,
	)

	requireNoError(t, err)
	if req.Resource != key.Resource {
		t.Fatalf("resource = %#v; want %#v", req.Resource, key.Resource)
	}
	if req.Object != key.Object {
		t.Fatalf("object = %#v; want %#v", req.Object, key.Object)
	}
	if req.Expected != 7 {
		t.Fatalf("expected = %s; want 7", req.Expected)
	}
	if req.Owner != owner {
		t.Fatalf("owner = %q; want %q", req.Owner, owner)
	}
	if req.Labels["app"] == nil || *req.Labels["app"] != "api" {
		t.Fatalf("label app = %#v; want api", req.Labels["app"])
	}
	if req.Annotations["note"] == nil || *req.Annotations["note"] != "owned" {
		t.Fatalf("annotation note = %#v; want owned", req.Annotations["note"])
	}
}

func TestPatchMetadataClonesMapsAndPointedValues(t *testing.T) {
	current := newCurrent(t, testKey("task-1"), 7, "desired")
	label := "api"
	annotation := "owned"
	labels := map[string]*string{"app": &label, "remove": nil}
	annotations := map[string]*string{"note": &annotation, "drop": nil}

	req, err := current.PatchMetadata(labels, annotations, testOwner())
	requireNoError(t, err)
	labels["extra"] = &label
	annotations["extra"] = &annotation
	label = "mutated"
	annotation = "mutated"

	if _, ok := req.Labels["extra"]; ok {
		t.Fatal("labels map was not detached")
	}
	if _, ok := req.Annotations["extra"]; ok {
		t.Fatal("annotations map was not detached")
	}
	if req.Labels["app"] == nil || *req.Labels["app"] != "api" {
		t.Fatalf("label pointer was not cloned: %#v", req.Labels["app"])
	}
	if req.Annotations["note"] == nil || *req.Annotations["note"] != "owned" {
		t.Fatalf("annotation pointer was not cloned: %#v", req.Annotations["note"])
	}
	if req.Labels["remove"] != nil {
		t.Fatalf("nil label patch became %#v", req.Labels["remove"])
	}
	if req.Annotations["drop"] != nil {
		t.Fatalf("nil annotation patch became %#v", req.Annotations["drop"])
	}
}

func TestPatchMetadataPreservesNilMaps(t *testing.T) {
	current := newCurrent(t, testKey("task-1"), 7, "desired")

	req, err := current.PatchMetadata(nil, nil, testOwner())

	requireNoError(t, err)
	if req.Labels != nil {
		t.Fatalf("labels = %#v; want nil", req.Labels)
	}
	if req.Annotations != nil {
		t.Fatalf("annotations = %#v; want nil", req.Annotations)
	}
}

func TestPatchMetadataRejectsZeroCurrent(t *testing.T) {
	var current Current

	_, err := current.PatchMetadata(nil, nil, testOwner())

	requireErrorIs(t, err, ErrInvalidCurrent)
}
