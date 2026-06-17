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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestErrorBuildersExposeSentinelsAndStructuredDiagnostics(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		reason   ErrorReason
		pathPart string
	}{
		{
			name:     "snapshot",
			err:      invalidSnapshotError("objectstorewatch.snapshot", objectstore.ErrInvalidListResult),
			target:   ErrInvalidSnapshot,
			reason:   ErrorReasonInvalidSnapshot,
			pathPart: "snapshot",
		},
		{
			name:     "boundary",
			err:      invalidBoundaryError("objectstorewatch.boundary", objectstore.ErrInvalidListRequest),
			target:   ErrInvalidBoundary,
			reason:   ErrorReasonInvalidBoundary,
			pathPart: "boundary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, tt.target)
			requireWatchError(t, tt.err, tt.reason, tt.pathPart)
		})
	}
}

func TestObjectStoreWatchErrorPreservesAllCauses(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")

	err := objectStoreWatchError(
		"objectstorewatch.snapshot",
		ErrorReasonInvalidSnapshot,
		[]error{ErrInvalidSnapshot},
		first,
		nil,
		second,
	)

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, first)
	requireErrorIs(t, err, second)
	requireWatchError(t, err, ErrorReasonInvalidSnapshot, "snapshot")
}

func TestObjectStoreWatchErrorUsesSentinelAsFallbackCause(t *testing.T) {
	err := objectStoreWatchError(
		"objectstorewatch.boundary",
		ErrorReasonInvalidBoundary,
		[]error{ErrInvalidBoundary},
	)

	requireErrorIs(t, err, ErrInvalidBoundary)
	requireWatchError(t, err, ErrorReasonInvalidBoundary, "boundary")
}

func TestJoinedCausesIgnoresNilValues(t *testing.T) {
	cause := errors.New("cause")

	joined := joinedCauses(nil, cause, nil)

	requireErrorIs(t, joined, cause)
}
