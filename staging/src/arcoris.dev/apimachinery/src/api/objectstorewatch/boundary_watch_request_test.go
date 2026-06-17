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

package objectstorewatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestBoundaryWatchRequestCreatesMatchingObjectWatchRequest(t *testing.T) {
	boundary, err := NewBoundary(testCollection(), 12)
	requireNoError(t, err)

	request, err := boundary.WatchRequest(WatchOptions{AllowProgress: true})
	requireNoError(t, err)

	if request.Collection != boundary.Collection() {
		t.Fatalf("Collection = %#v; want %#v", request.Collection, boundary.Collection())
	}
	if request.Start.Mode != objectwatch.StartAfterRevision {
		t.Fatalf("Start.Mode = %s; want %s", request.Start.Mode, objectwatch.StartAfterRevision)
	}
	if request.Start.Revision != boundary.Revision() {
		t.Fatalf("Start.Revision = %s; want %s", request.Start.Revision, boundary.Revision())
	}
	if !request.AllowProgress {
		t.Fatalf("AllowProgress = false; want true")
	}
}

func TestBoundaryWatchRequestPreservesZeroRevisionBoundary(t *testing.T) {
	boundary, err := NewBoundary(testCollection(), 0)
	requireNoError(t, err)

	request, err := boundary.WatchRequest(WatchOptions{})
	requireNoError(t, err)

	if request.Start.Mode != objectwatch.StartAfterRevision {
		t.Fatalf("Start.Mode = %s; want %s", request.Start.Mode, objectwatch.StartAfterRevision)
	}
	if !request.Start.Revision.IsZero() {
		t.Fatalf("Start.Revision = %s; want zero", request.Start.Revision)
	}
	if request.AllowProgress {
		t.Fatalf("AllowProgress = true; want false")
	}
}

func TestBoundaryWatchRequestRejectsInvalidBoundary(t *testing.T) {
	request, err := (Boundary{}).WatchRequest(WatchOptions{AllowProgress: true})

	requireErrorIs(t, err, ErrInvalidBoundary)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
	if request != (objectwatch.Request{}) {
		t.Fatalf("request = %#v; want zero request", request)
	}
}
