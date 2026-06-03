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

package objectapply

import (
	"testing"
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/finalizer"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	metaowner "arcoris.dev/apimachinery/api/meta/owner"
	"arcoris.dev/apimachinery/api/meta/stamp"
)

func TestApplyPreservesLiveTypeMeta(t *testing.T) {
	result, err := Apply(testRequest(), Options{})
	requireNoError(t, err)

	if result.Object.TypeMeta != testTypeMeta("v1") {
		t.Fatalf("TypeMeta = %#v; want live TypeMeta", result.Object.TypeMeta)
	}
}

func TestApplyPreservesLiveObjectMeta(t *testing.T) {
	result, err := Apply(testRequest(), Options{})
	requireNoError(t, err)

	if result.Object.ObjectMeta.ResourceVersion != testObjectMeta().ResourceVersion {
		t.Fatalf("ResourceVersion = %q; want live metadata", result.Object.ObjectMeta.ResourceVersion)
	}
	if result.Object.ObjectMeta.Generation != testObjectMeta().Generation {
		t.Fatalf("Generation = %d; want live metadata", result.Object.ObjectMeta.Generation)
	}
}

func TestApplyDoesNotCopyAppliedTypeMeta(t *testing.T) {
	req := testRequest()

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	if result.Object.TypeMeta != req.Live.TypeMeta {
		t.Fatalf("TypeMeta = %#v; want live TypeMeta", result.Object.TypeMeta)
	}
}

func TestApplyDoesNotCopyAppliedObjectMeta(t *testing.T) {
	req := testRequest()
	req.Live.ObjectMeta.ResourceVersion = stamp.ResourceVersion("rv-live")

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	if result.Object.ObjectMeta.ResourceVersion != "rv-live" {
		t.Fatalf("ResourceVersion = %q; want live metadata", result.Object.ObjectMeta.ResourceVersion)
	}
}

func TestApplyRejectsUnsupportedMetadataChange(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Labels = labels.Set{"role": "worker"}

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrUnsupportedMetadataChange)
}

func TestApplyRejectsUnsupportedAppliedMetadataFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*ValueObject)
	}{
		{
			name: "generateName",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.GenerateName = metaidentity.NamePrefix("worker-")
			},
		},
		{
			name: "resourceVersion",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.ResourceVersion = stamp.ResourceVersion("rv-applied")
			},
		},
		{
			name: "generation",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.Generation = stamp.Generation(8)
			},
		},
		{
			name: "createdAt",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.CreatedAt = metadataTimestamp()
			},
		},
		{
			name: "deletion",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.Deletion = &stamp.Deletion{DeletedAt: metadataTimestamp()}
			},
		},
		{
			name: "labels",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.Labels = labels.Set{"role": "worker"}
			},
		},
		{
			name: "annotations",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.Annotations = annotations.Set{"control.arcoris.dev/note": "worker"}
			},
		},
		{
			name: "ownerReferences",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.OwnerReferences = metaowner.List{metadataOwnerReference()}
			},
		},
		{
			name: "finalizers",
			mutate: func(obj *ValueObject) {
				obj.ObjectMeta.Finalizers = finalizer.Set{finalizer.Name("cleanup")}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testRequest()
			tt.mutate(&req.Applied)

			_, err := Apply(req, Options{})

			requireErrorIs(t, err, ErrUnsupportedMetadataChange)
			requireObjectApplyError(
				t,
				err,
				pathObjectAppliedMetadata,
				ErrorReasonUnsupportedMetadataChange,
			)
		})
	}
}

func TestApplyAllowsEquivalentIdentityMetadata(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta = minimalAppliedObjectMeta()

	_, err := Apply(req, Options{})

	requireNoError(t, err)
}

func TestApplyAllowsNonNilZeroAppliedDeletion(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Deletion = &stamp.Deletion{}

	_, err := Apply(req, Options{})

	requireNoError(t, err)
}

func TestApplyRejectsNonZeroAppliedDeletion(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Deletion = &stamp.Deletion{DeletedAt: metadataTimestamp()}

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrUnsupportedMetadataChange)
}

func metadataTimestamp() stamp.Timestamp {
	return stamp.NewTimestamp(time.Date(2026, 6, 3, 12, 0, 0, 0, time.UTC))
}

func metadataOwnerReference() metaowner.Reference {
	return metaowner.Reference{
		Ref: metaidentity.ObjectReference{
			APIVersion: apiidentity.GroupVersion{
				Group:   "control.arcoris.dev",
				Version: "v1",
			},
			Kind:      "Worker",
			Namespace: "system",
			Name:      "parent",
			UID:       "uid-parent",
		},
	}
}
