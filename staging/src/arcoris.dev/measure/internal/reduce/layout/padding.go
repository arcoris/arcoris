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

package layout

import "arcoris.dev/atomicx"

const (
	// CacheLinePadSize is the canonical ARCORIS cache-line padding width.
	CacheLinePadSize = atomicx.CacheLinePadSize
)

// CacheLinePad aliases the canonical ARCORIS cache-line pad type.
//
// Keeping the alias here lets reduction code share the same cache-line contract
// as atomicx without copying architecture assumptions into this package.
type CacheLinePad = atomicx.CacheLinePad

// Padded stores Value between cache-line pads.
//
// Padded is a layout utility for rare cases where a caller intentionally places
// independently hot worker slots in adjacent memory. It is not used by default
// runners because workers compute into local variables and write one partial
// slot once.
type Padded[T any] struct {
	// _leading separates Value from fields or slice elements that may precede this
	// slot in memory.
	_leading CacheLinePad

	// Value is the caller-owned payload protected from adjacent false sharing by
	// the leading and trailing pads.
	Value T

	// _trailing separates Value from fields or slice elements that may follow this
	// slot in memory.
	_trailing CacheLinePad
}
