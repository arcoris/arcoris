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

	"arcoris.dev/apimachinery/api/meta/labels"
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

func TestApplyAllowsEquivalentIdentityMetadata(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta = minimalAppliedObjectMeta()

	_, err := Apply(req, Options{})

	requireNoError(t, err)
}
