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

import "context"

// Stream is a pull-based stream of validated object watch events.
//
// Next blocks until an event, context cancellation, stream closure, or a
// terminal error. A stream must not end with a normal source-side EOF: source
// termination without caller cancellation, Close, EventRestartRequired,
// ErrHistoryUnavailable, or ErrContinuityLost is itself a continuity condition.
// Implementations must not emit events after terminal restart/continuity loss.
type Stream interface {
	// Next returns the next validated event or a terminal/context error.
	Next(context.Context) (Event, error)
	// Close releases stream resources. Close must be idempotent.
	Close() error
}
