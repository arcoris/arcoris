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

// LookupValue returns a byte-based value codec for mediaType.
func (r Registry) LookupValue(mediaType codec.MediaType) (codec.ValueCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	valueCodec, ok := entry.codec.(codec.ValueCodec)
	return valueCodec, ok
}

// LookupObject returns a byte-based object codec for mediaType.
func (r Registry) LookupObject(mediaType codec.MediaType) (codec.ObjectCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	objectCodec, ok := entry.codec.(codec.ObjectCodec)
	return objectCodec, ok
}

// LookupObjectOwnership returns a byte-based ownership codec for mediaType.
func (r Registry) LookupObjectOwnership(mediaType codec.MediaType) (codec.ObjectOwnershipCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	ownershipCodec, ok := entry.codec.(codec.ObjectOwnershipCodec)
	return ownershipCodec, ok
}

// LookupCodec returns a full byte-based codec for mediaType.
func (r Registry) LookupCodec(mediaType codec.MediaType) (codec.Codec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	fullCodec, ok := entry.codec.(codec.Codec)
	return fullCodec, ok
}
