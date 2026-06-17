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
	"fmt"
)

// ErrorReason identifies a stable objectwatch diagnostic class.
type ErrorReason string

const (
	// ErrorReasonInvalidRequest identifies a malformed watch request.
	ErrorReasonInvalidRequest ErrorReason = "invalid_request"
	// ErrorReasonInvalidStart identifies a malformed or unsupported start.
	ErrorReasonInvalidStart ErrorReason = "invalid_start"
	// ErrorReasonInvalidEvent identifies a malformed stream event.
	ErrorReasonInvalidEvent ErrorReason = "invalid_event"
	// ErrorReasonInvalidRestart identifies a malformed restart reason.
	ErrorReasonInvalidRestart ErrorReason = "invalid_restart"
	// ErrorReasonClosed identifies use after a terminal event or close.
	ErrorReasonClosed ErrorReason = "closed"
	// ErrorReasonHistoryUnavailable identifies unavailable requested history.
	ErrorReasonHistoryUnavailable ErrorReason = "history_unavailable"
	// ErrorReasonContinuityLost identifies an observed stream ordering gap.
	ErrorReasonContinuityLost ErrorReason = "continuity_lost"
)

// Error carries a structured objectwatch diagnostic while preserving causes.
type Error struct {
	// Path identifies the logical contract location that failed.
	Path string
	// Reason is the stable machine-readable failure class.
	Reason ErrorReason
	// Cause is the lower validation or stream failure.
	Cause error
}

// Error returns compact diagnostic text.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Path == "" {
		return string(e.Reason) + ": " + e.Cause.Error()
	}

	return e.Path + ": " + string(e.Reason) + ": " + e.Cause.Error()
}

// Unwrap returns the lower cause for errors.Is and errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

var (
	// ErrInvalidRequest classifies malformed object watch requests.
	ErrInvalidRequest = errors.New("invalid object watch request")
	// ErrInvalidStart classifies malformed or unsupported watch start points.
	ErrInvalidStart = errors.New("invalid object watch start")
	// ErrInvalidEvent classifies malformed object watch events.
	ErrInvalidEvent = errors.New("invalid object watch event")
	// ErrInvalidRestart classifies malformed restart-required reasons.
	ErrInvalidRestart = errors.New("invalid object watch restart")
	// ErrClosed classifies use after a terminal stream state.
	ErrClosed = errors.New("object watch closed")
	// ErrHistoryUnavailable reports that requested history cannot be served.
	ErrHistoryUnavailable = errors.New("object watch history unavailable")
	// ErrContinuityLost reports that a stream can no longer prove continuity.
	ErrContinuityLost = errors.New("object watch continuity lost")
)

// invalidStartError builds an error rooted at watch.start.
func invalidStartError(format string, args ...any) error {
	return objectWatchError(
		"watch.start",
		ErrorReasonInvalidStart,
		[]error{ErrInvalidStart},
		fmt.Errorf(format, args...),
	)
}

// invalidRestartError builds an error rooted at watch.event.restart.
func invalidRestartError(cause error) error {
	return objectWatchError(
		"watch.event.restart",
		ErrorReasonInvalidRestart,
		[]error{ErrInvalidRestart},
		cause,
	)
}

// invalidEventError builds an event validation error at path.
func invalidEventError(path string, cause error, sentinels ...error) error {
	return objectWatchError(
		path,
		ErrorReasonInvalidEvent,
		append([]error{ErrInvalidEvent}, sentinels...),
		cause,
	)
}

// invalidRequestError builds a request validation error at path.
func invalidRequestError(path string, cause error, sentinels ...error) error {
	return objectWatchError(
		path,
		ErrorReasonInvalidRequest,
		append([]error{ErrInvalidRequest}, sentinels...),
		cause,
	)
}

// continuityLostError reports an observed event sequence violation.
func continuityLostError(cause error) error {
	return objectWatchError(
		"watch.validator",
		ErrorReasonContinuityLost,
		[]error{ErrContinuityLost},
		cause,
	)
}

// closedError reports use after a terminal event.
func closedError(cause error) error {
	return objectWatchError(
		"watch.validator",
		ErrorReasonClosed,
		[]error{ErrClosed},
		cause,
	)
}

// objectWatchError joins broad sentinels, lower causes, and one structured
// diagnostic without hiding errors.Is/errors.As behavior.
func objectWatchError(path string, reason ErrorReason, sentinels []error, causes ...error) error {
	var cause error
	for _, item := range causes {
		if item == nil {
			continue
		}
		if cause == nil {
			cause = item
		} else {
			cause = errors.Join(cause, item)
		}
	}
	if cause == nil && len(sentinels) > 0 {
		cause = sentinels[len(sentinels)-1]
	}
	if cause == nil {
		cause = errors.New("object watch contract violation")
	}

	diagnostic := &Error{Path: path, Reason: reason, Cause: cause}
	joined := make([]error, 0, len(sentinels)+len(causes)+1)
	joined = append(joined, sentinels...)
	joined = append(joined, diagnostic)
	for _, item := range causes {
		if item != nil {
			joined = append(joined, item)
		}
	}
	return errors.Join(joined...)
}
