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

package listmapkey

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestIsPayloadFailure(t *testing.T) {
	err := failure(fieldpath.RootPath(), FailureMissingKey, "missing")

	requireEqual(t, IsPayloadFailure(err), true)
	requireEqual(t, IsDescriptorFailure(err), false)
}

func TestIsDescriptorFailure(t *testing.T) {
	err := failure(fieldpath.RootPath(), FailureUnresolvedRef, "unresolved")

	requireEqual(t, IsDescriptorFailure(err), true)
	requireEqual(t, IsPayloadFailure(err), false)
}

func TestFailureClassificationCoversAllKinds(t *testing.T) {
	tests := []struct {
		name       string
		kind       FailureKind
		payload    bool
		descriptor bool
	}{
		{name: "invalid descriptor", kind: FailureInvalidDescriptor, descriptor: true},
		{name: "unresolved ref", kind: FailureUnresolvedRef, descriptor: true},
		{name: "reference cycle", kind: FailureReferenceCycle, descriptor: true},
		{name: "item kind mismatch", kind: FailureItemKindMismatch, payload: true},
		{name: "missing key", kind: FailureMissingKey, payload: true},
		{name: "null key", kind: FailureNullKey, payload: true},
		{name: "key kind mismatch", kind: FailureKeyKindMismatch, payload: true},
		{name: "key integer range", kind: FailureKeyIntegerRange, payload: true},
		{name: "invalid selector", kind: FailureInvalidSelector, descriptor: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := failure(fieldpath.RootPath(), tt.kind, "detail")

			requireEqual(t, IsPayloadFailure(err), tt.payload)
			requireEqual(t, IsDescriptorFailure(err), tt.descriptor)
			requireEqual(t, IsPayloadFailure(err) && IsDescriptorFailure(err), false)
		})
	}
}

func TestFailureClassificationIgnoresUnrelatedErrors(t *testing.T) {
	err := errors.New("not a listmapkey error")

	requireEqual(t, IsPayloadFailure(err), false)
	requireEqual(t, IsDescriptorFailure(err), false)
}
