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
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectownership"
)

func TestApplySameDesiredIntentIsObjectAndOwnershipIdempotent(t *testing.T) {
	req := testRequest()
	req.Live = testObjectObserved(req.Live.Desired, obj(member("ready", str("true"))))
	req.Resource = testResourceWithObserved(desiredDescriptor())
	req.Ownership = objectownership.NewStateWithSurfaces(
		fieldownership.MustState(entry("other", path("$.replicas"))),
		fieldownership.MustState(entry("controller", path("$.ready"))),
		objectownership.NewMetadataState(
			fieldownership.MustState(entry("labeler", path(`$["app"]`))),
			fieldownership.MustState(entry("annotator", path(`$["note"]`))),
		),
	)

	first, err := Apply(req, Options{})
	requireNoError(t, err)

	req.Live = first.Object
	req.Ownership = first.Ownership
	second, err := Apply(req, Options{})
	requireNoError(t, err)

	if !reflect.DeepEqual(second.Object, first.Object) {
		t.Fatalf("second object = %#v; want %#v", second.Object, first.Object)
	}
	if !reflect.DeepEqual(second.Ownership, first.Ownership) {
		t.Fatalf("second ownership = %#v; want %#v", second.Ownership, first.Ownership)
	}

	requireStringMember(t, second.Object.Desired, "image", "api:v2")
	requireStringMember(t, *second.Object.Observed, "ready", "true")
	requireOwnersOf(t, second.Ownership.Desired(), path("$.image"), "user")
	requireOwnersOf(t, second.Ownership.Desired(), path("$.replicas"), "other")
	requireOwnersOf(t, second.Ownership.Observed(), path("$.ready"), "controller")
	requireOwnersOf(t, second.Ownership.Metadata().Labels(), path(`$["app"]`), "labeler")
	requireOwnersOf(t, second.Ownership.Metadata().Annotations(), path(`$["note"]`), "annotator")
}
