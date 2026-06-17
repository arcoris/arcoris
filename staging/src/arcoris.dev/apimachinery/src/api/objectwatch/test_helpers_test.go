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

package objectwatch

import (
	"errors"
	"strings"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// watchResource returns the stable resource collection used by watch tests.
func watchResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

// watchListRequest returns a structurally valid all-namespaces watch request
// collection.
func watchListRequest() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: watchResource(),
		Scope:    objectstore.AllNamespaces(),
	}
}

// watchKey returns an authoritative storage key for the shared test resource.
func watchKey(name metaidentity.Name) objectstore.Key {
	return objectstore.MustKey(watchResource(), metaidentity.ObjectName{
		Namespace: "system",
		Name:      name,
	})
}

// watchState returns committed object state with a chosen revision.
func watchState(revision objectstore.Revision, desired string) objectstore.State {
	return objectstore.State{
		Object: object.New[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{
				Name:      "main",
				Namespace: "system",
			},
			value.StringValue(desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}

// watchCreatedChange returns a valid created transition fixture.
func watchCreatedChange(revision objectstore.Revision) objectstore.Change {
	return objectstore.MustCreatedChange(watchKey("main"), watchState(revision, "created"))
}

// watchUpdatedChange returns a valid updated transition fixture.
func watchUpdatedChange(before objectstore.Revision, after objectstore.Revision) objectstore.Change {
	return objectstore.MustUpdatedChange(
		watchKey("main"),
		watchState(before, "before"),
		watchState(after, "after"),
	)
}

// watchDeletedChange returns a valid deleted transition fixture.
func watchDeletedChange(before objectstore.Revision, revision objectstore.Revision) objectstore.Change {
	return objectstore.MustDeletedChange(watchKey("main"), watchState(before, "before"), revision)
}

// requireNoError fails the current test when err is non-nil.
func requireNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs fails unless err preserves target through errors.Is.
func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

// requireWatchError verifies the structured objectwatch diagnostic selected by
// errors.As without depending on the full Error string.
func requireWatchError(t testing.TB, err error, reason ErrorReason, pathPart string) {
	t.Helper()

	var watchErr *Error
	if !errors.As(err, &watchErr) {
		t.Fatalf("errors.As(%v, *Error) = false", err)
	}
	if watchErr.Reason != reason {
		t.Fatalf("reason = %s; want %s", watchErr.Reason, reason)
	}
	if !strings.Contains(watchErr.Path, pathPart) {
		t.Fatalf("path = %q; want to contain %q", watchErr.Path, pathPart)
	}
}
