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

package objectstore

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func TestValidateInputStateAcceptsZeroRevision(t *testing.T) {
	requireNoError(t, ValidateInputState(validState()))
}

func TestValidateInputStateAcceptsOwnership(t *testing.T) {
	state := validState()
	state.Ownership = ownershipWithSurfaces()

	requireNoError(t, ValidateInputState(state))
}

func TestValidateInputStateRejectsForgedRevision(t *testing.T) {
	state := validCommittedState()

	requireErrorIs(t, ValidateInputState(state), ErrInvalidRevision)
}

func TestValidateCommittedStateRequiresRevision(t *testing.T) {
	requireErrorIs(t, ValidateCommittedState(validState()), ErrInvalidRevision)
}

func TestValidateCommittedStateAcceptsCommittedState(t *testing.T) {
	requireNoError(t, ValidateCommittedState(validCommittedState()))
}

func TestValidateCommittedStateRejectsInvalidObjectMeta(t *testing.T) {
	state := validCommittedState()
	state.Object.Kind = ""

	requireErrorIs(t, ValidateCommittedState(state), ErrInvalidState)
	requireErrorIs(t, ValidateCommittedState(state), object.ErrInvalidObject)
	requireStoreErrorReason(t, ValidateCommittedState(state), ErrorReasonInvalidStateObject)
}

func TestValidateCommittedStateRejectsInvalidDesiredValue(t *testing.T) {
	state := validCommittedState()
	state.Object.Desired = value.Value{}

	requireErrorIs(t, ValidateCommittedState(state), ErrInvalidState)
	requireStoreErrorReason(t, ValidateCommittedState(state), ErrorReasonMissingDesired)
}

func TestValidateCommittedStateRejectsInvalidObservedValue(t *testing.T) {
	state := validCommittedState()
	state.Object.Observed = &value.Value{}

	requireErrorIs(t, ValidateCommittedState(state), ErrInvalidState)
	requireStoreErrorReason(t, ValidateCommittedState(state), ErrorReasonInvalidObserved)
}

func TestValidateCommittedStateAcceptsNormalizedOwnership(t *testing.T) {
	state := validCommittedState()
	state.Ownership = objectownership.Normalize(ownershipWithSurfaces())

	requireNoError(t, ValidateCommittedState(state))
}

func requireStoreErrorReason(t *testing.T, err error, reason ErrorReason) {
	t.Helper()

	var storeErr *Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("errors.As(%v, *Error) = false", err)
	}
	if storeErr.Reason != reason {
		t.Fatalf("reason = %v; want %v", storeErr.Reason, reason)
	}
}
