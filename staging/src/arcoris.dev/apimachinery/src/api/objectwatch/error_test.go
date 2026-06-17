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
