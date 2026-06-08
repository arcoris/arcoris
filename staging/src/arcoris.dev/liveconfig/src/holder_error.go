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

package liveconfig

import "errors"

// ErrNilHolder reports a method call on a nil or unconstructed Holder receiver.
//
// Holder methods panic with this value instead of returning zero snapshots from
// a nil or zero-value receiver. Such a Holder has no publisher, no last-good
// value, and no meaningful revision, so treating it as a programming error
// keeps failures explicit.
var ErrNilHolder = errors.New("liveconfig: nil holder")

// requireHolder panics when h is nil or was not constructed with New.
func requireHolder[T any](h *Holder[T]) {
	if h == nil || h.pub == nil {
		panic(ErrNilHolder)
	}
}
