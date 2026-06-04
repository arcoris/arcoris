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

package codec

// DecodeOptions controls common format-independent decode behavior.
//
// Zero DecodeOptions are valid. Concrete codec packages may define additional
// format-specific options around this common intent. These options describe
// decoding policy only; they do not request descriptor validation, defaulting,
// pruning, conversion, resource lookup, or apply behavior.
type DecodeOptions struct {
	// Strict requests safe, unambiguous decoding.
	//
	// Concrete codecs define exact strict behavior for their format. For
	// structured document formats, strict mode should reject ambiguous or lossy
	// input when detectable.
	Strict bool

	// MaxDepth bounds nested input structures. Zero means implementation default.
	//
	// Concrete codecs choose the package default and the exact definition of
	// nesting depth for their format.
	MaxDepth int

	// RejectDuplicateFields rejects duplicate object/map keys when the format
	// and implementation can detect them. Strict mode should imply this behavior
	// for structured formats with object/map keys.
	RejectDuplicateFields bool
}
