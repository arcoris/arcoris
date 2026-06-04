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

package codecregistry

import "arcoris.dev/apimachinery/api/codec"

// Entry is one normalized codec registration record.
//
// Entry stores the original already-configured codec implementation together
// with normalized, detached metadata accepted by Registry construction. It does
// not call Info again after registration.
type Entry struct {
	// codec is the original implementation passed to New.
	//
	// The registry keeps this exact interface value so typed lookup methods can
	// return the implementation without wrappers or adapter allocations.
	codec codec.BaseCodec

	// info is the normalized metadata accepted during construction.
	//
	// Keeping a detached snapshot makes registry behavior independent from any
	// later changes in a codec implementation's Info method or backing slices.
	info codec.Info
}

// IsZero reports whether e contains no codec and no metadata.
func (e Entry) IsZero() bool {
	return e.codec == nil && e.info.IsZero()
}

// Codec returns the original codec implementation.
func (e Entry) Codec() codec.BaseCodec {
	return e.codec
}

// Info returns detached normalized metadata for the entry.
func (e Entry) Info() codec.Info {
	return cloneInfo(e.info)
}

// Format returns the entry's normalized format.
func (e Entry) Format() codec.Format {
	return e.info.Format
}

// MediaTypes returns detached normalized media types for the entry.
func (e Entry) MediaTypes() []codec.MediaType {
	return e.Info().MediaTypes
}

// Targets returns detached normalized semantic targets for the entry.
func (e Entry) Targets() []codec.Target {
	return e.Info().Targets
}
