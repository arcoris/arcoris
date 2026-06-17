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

// Source opens watch streams for structural object collections.
//
// Source implementations own the concrete continuity proof. Watch must
// validate Request before opening a stream. Malformed requests should return
// ErrInvalidRequest and/or ErrInvalidStart. Valid requests unsupported by the
// source should return ErrUnsupportedCapability. If requested history cannot be
// served, Watch should return ErrHistoryUnavailable. If continuity is lost
// after stream creation, the stream should emit EventRestartRequired or return
// ErrContinuityLost.
type Source interface {
	// Watch opens a pull-based event stream for request.
	//
	// If Watch returns nil error, the Stream must be non-nil. If Watch returns
	// a non-nil error, the Stream must be nil.
	Watch(context.Context, Request) (Stream, error)
}
