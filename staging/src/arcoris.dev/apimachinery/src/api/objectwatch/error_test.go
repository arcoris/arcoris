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
	"testing"
)

func TestErrorsExposeSentinelsAndStructuredDiagnostics(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		reason   ErrorReason
		pathPart string
	}{
		{name: "invalid request", err: invalidRequestError("watch.request", ErrInvalidRequest), target: ErrInvalidRequest, reason: ErrorReasonInvalidRequest, pathPart: "watch.request"},
		{name: "invalid start", err: invalidStartError("bad"), target: ErrInvalidStart, reason: ErrorReasonInvalidStart, pathPart: "watch.start"},
		{name: "unsupported capability", err: UnsupportedCapability(nil), target: ErrUnsupportedCapability, reason: ErrorReasonUnsupportedCapability, pathPart: "watch.capabilities"},
		{name: "invalid event", err: invalidEventError("watch.event", ErrInvalidEvent), target: ErrInvalidEvent, reason: ErrorReasonInvalidEvent, pathPart: "watch.event"},
		{name: "invalid restart", err: invalidRestartError(ErrInvalidRestart), target: ErrInvalidRestart, reason: ErrorReasonInvalidRestart, pathPart: "watch.event.restart"},
		{name: "closed", err: closedError(ErrClosed), target: ErrClosed, reason: ErrorReasonClosed, pathPart: "watch.validator"},
		{name: "history unavailable", err: objectWatchError("watch.source", ErrorReasonHistoryUnavailable, []error{ErrHistoryUnavailable}, ErrHistoryUnavailable), target: ErrHistoryUnavailable, reason: ErrorReasonHistoryUnavailable, pathPart: "watch.source"},
		{name: "continuity lost", err: continuityLostError(ErrContinuityLost), target: ErrContinuityLost, reason: ErrorReasonContinuityLost, pathPart: "watch.validator"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, tt.target)
			requireWatchError(t, tt.err, tt.reason, tt.pathPart)
		})
	}
}

func TestErrorUnwrapPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := invalidRequestError("watch.request", cause)

	requireErrorIs(t, err, cause)
}

func TestExportedErrorConstructorsPreserveCause(t *testing.T) {
	cause := errors.New("source failed")
	tests := []struct {
		name     string
		err      error
		target   error
		reason   ErrorReason
		pathPart string
	}{
		{name: "history unavailable", err: HistoryUnavailable(cause), target: ErrHistoryUnavailable, reason: ErrorReasonHistoryUnavailable, pathPart: "watch.source.history"},
		{name: "continuity lost", err: ContinuityLost(cause), target: ErrContinuityLost, reason: ErrorReasonContinuityLost, pathPart: "watch.source.continuity"},
		{name: "closed", err: Closed(cause), target: ErrClosed, reason: ErrorReasonClosed, pathPart: "watch.stream.closed"},
		{name: "unsupported capability", err: UnsupportedCapability(cause), target: ErrUnsupportedCapability, reason: ErrorReasonUnsupportedCapability, pathPart: "watch.capabilities"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, tt.target)
			requireErrorIs(t, tt.err, cause)
			requireWatchError(t, tt.err, tt.reason, tt.pathPart)
		})
	}
}

func TestExportedErrorConstructorsAcceptNilCause(t *testing.T) {
	tests := []error{
		HistoryUnavailable(nil),
		ContinuityLost(nil),
		Closed(nil),
		UnsupportedCapability(nil),
	}

	for _, err := range tests {
		if err == nil {
			t.Fatalf("constructor returned nil")
		}
		var watchErr *Error
		if !errors.As(err, &watchErr) {
			t.Fatalf("errors.As(%v, *Error) = false", err)
		}
	}
}

func TestErrorStringHandlesNilCause(t *testing.T) {
	tests := []*Error{
		nil,
		&Error{},
		&Error{Reason: ErrorReasonInvalidEvent},
		&Error{Path: "watch.event", Reason: ErrorReasonInvalidEvent},
	}

	for _, err := range tests {
		if got := err.Error(); got == "" {
			t.Fatalf("Error() returned empty string")
		}
	}
}
