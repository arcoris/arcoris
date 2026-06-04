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

// Info describes a codec implementation.
//
// It contains only stable metadata for concrete codecs and future registries.
// It does not contain implementation pointers, registry indexes, priority
// rules, or process-wide registration state. Codec implementations should
// return detached metadata slices from Info() so callers cannot mutate shared
// implementation state.
type Info struct {
	// Format identifies the codec's format family.
	//
	// The format is open-world but must be a canonical Format token when the
	// Info is validated.
	Format Format

	// MediaTypes lists concrete media types accepted or emitted by the codec.
	//
	// The list must be non-empty, duplicate-free, and canonical after
	// validation or normalization.
	MediaTypes []MediaType

	// Targets lists API document models supported by the codec.
	//
	// Targets describe semantic API document models, not transport
	// capabilities. Byte-slice support is represented by ValueCodec,
	// ObjectCodec, and ObjectOwnershipCodec. Streaming support is represented
	// by ValueStreamCodec, ObjectStreamCodec, and ObjectOwnershipStreamCodec.
	// The list must be non-empty and contain only the closed-world Target values
	// known to this package.
	Targets []Target
}

// IsZero reports whether i contains no codec metadata.
//
// IsZero is an absence check only. It does not imply that a non-zero Info is
// valid, canonical, sorted, or duplicate-free.
func (i Info) IsZero() bool {
	return i.Format.IsZero() && len(i.MediaTypes) == 0 && len(i.Targets) == 0
}
