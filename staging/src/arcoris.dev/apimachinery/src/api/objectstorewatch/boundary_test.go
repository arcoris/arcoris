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
)

func TestNewBoundaryAcceptsZeroRevision(t *testing.T) {
	boundary, err := NewBoundary(testCollection(), 0)
	requireNoError(t, err)

	if !boundary.Revision().IsZero() {
		t.Fatalf("Revision() = %s; want zero", boundary.Revision())
	}
	if boundary.Collection() != testCollection() {
		t.Fatalf("Collection() = %#v; want %#v", boundary.Collection(), testCollection())
	}
}

func TestNewBoundaryAcceptsNonZeroRevision(t *testing.T) {
	boundary, err := NewBoundary(testCollection(), 12)
	requireNoError(t, err)

	if boundary.Revision() != 12 {
		t.Fatalf("Revision() = %s; want 12", boundary.Revision())
	}
}

func TestNewBoundaryRejectsInvalidCollection(t *testing.T) {
	_, err := NewBoundary(objectstore.ListRequest{}, 0)

	requireErrorIs(t, err, ErrInvalidBoundary)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
	requireWatchError(t, err, ErrorReasonInvalidBoundary, "boundary.collection")
}

func TestBoundaryCloneReturnsEquivalentBoundary(t *testing.T) {
	boundary, err := NewBoundary(testCollection(), 7)
	requireNoError(t, err)

	cloned := boundary.Clone()
	if cloned.Collection() != boundary.Collection() {
		t.Fatalf("clone collection = %#v; want %#v", cloned.Collection(), boundary.Collection())
	}
	if cloned.Revision() != boundary.Revision() {
		t.Fatalf("clone revision = %s; want %s", cloned.Revision(), boundary.Revision())
	}
}

func TestBoundaryIsZero(t *testing.T) {
	if !(Boundary{}).IsZero() {
		t.Fatalf("zero Boundary reported non-zero")
	}

	boundary, err := NewBoundary(testCollection(), 0)
	requireNoError(t, err)
	if boundary.IsZero() {
		t.Fatalf("valid zero-revision boundary reported zero")
	}
}

func TestBoundaryIsValid(t *testing.T) {
	boundary, err := NewBoundary(testCollection(), 0)
	requireNoError(t, err)

	if !boundary.IsValid() {
		t.Fatalf("valid boundary reported invalid")
	}
	if (Boundary{}).IsValid() {
		t.Fatalf("zero boundary reported valid")
	}
}
