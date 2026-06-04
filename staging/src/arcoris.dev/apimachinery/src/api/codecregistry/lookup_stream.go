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

// LookupValueStream returns a streaming value codec for mediaType.
func (r Registry) LookupValueStream(mediaType codec.MediaType) (codec.ValueStreamCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	valueCodec, ok := entry.codec.(codec.ValueStreamCodec)
	return valueCodec, ok
}

// LookupObjectStream returns a streaming object codec for mediaType.
func (r Registry) LookupObjectStream(mediaType codec.MediaType) (codec.ObjectStreamCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	objectCodec, ok := entry.codec.(codec.ObjectStreamCodec)
	return objectCodec, ok
}

// LookupObjectOwnershipStream returns a streaming ownership codec for mediaType.
func (r Registry) LookupObjectOwnershipStream(
	mediaType codec.MediaType,
) (codec.ObjectOwnershipStreamCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	ownershipCodec, ok := entry.codec.(codec.ObjectOwnershipStreamCodec)
	return ownershipCodec, ok
}

// LookupStreamingCodec returns a full streaming codec for mediaType.
func (r Registry) LookupStreamingCodec(mediaType codec.MediaType) (codec.StreamingCodec, bool) {
	entry, ok := r.LookupMediaType(mediaType)
	if !ok {
		return nil, false
	}

	streamingCodec, ok := entry.codec.(codec.StreamingCodec)
	return streamingCodec, ok
}
