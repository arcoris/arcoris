/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package retrybudget

// Kind identifies the retry-budget implementation that produced a snapshot.
//
// Kind is diagnostic metadata. The root package defines only kinds for
// implementations that exist or are intentionally part of the first contract
// surface. Future implementations should add new kinds when they are introduced.
type Kind uint8

const (
	// KindUnknown is the zero kind and is not valid for produced snapshots.
	KindUnknown Kind = iota

	// KindNoop identifies an implementation that never exhausts retry budget.
	KindNoop

	// KindFixedWindow identifies a fixed-window traffic-ratio implementation.
	KindFixedWindow
)

// String returns the stable diagnostic name for k.
func (k Kind) String() string {
	switch k {
	case KindNoop:
		return "noop"
	case KindFixedWindow:
		return "fixed_window"
	default:
		return "unknown"
	}
}

// IsValid reports whether k identifies a concrete retry-budget implementation.
func (k Kind) IsValid() bool {
	switch k {
	case KindNoop, KindFixedWindow:
		return true
	default:
		return false
	}
}
