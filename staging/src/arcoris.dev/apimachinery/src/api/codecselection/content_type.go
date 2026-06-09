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

package codecselection

import "arcoris.dev/apimachinery/api/codec"

// ContentType is normalized protocol-neutral content identity key material.
//
// ContentType stores a canonical codec.MediaType plus deterministic parameters.
// It is not a parser for raw Content-Descriptor headers.
type ContentType struct {
	// mediaType is the canonical codec media type.
	mediaType codec.MediaType

	// parameters are normalized deterministic key parameters.
	parameters Parameters
}

// NewContentType validates mediaType and parameters into a normalized content type.
func NewContentType(mediaType codec.MediaType, parameters ...Parameter) (ContentType, error) {
	contentType := ContentType{mediaType: mediaType, parameters: Parameters{items: parameters}}

	return normalizeContentTypeAt(pathContentType, contentType)
}

// MustContentType returns a normalized ContentType or panics when input is invalid.
func MustContentType(mediaType codec.MediaType, parameters ...Parameter) ContentType {
	contentType, err := NewContentType(mediaType, parameters...)
	if err != nil {
		panic(err)
	}

	return contentType
}

// IsZero reports whether c contains no media type and no parameters.
func (c ContentType) IsZero() bool {
	return c.mediaType.IsZero() && c.parameters.IsZero()
}

// MediaType returns the canonical media type.
func (c ContentType) MediaType() codec.MediaType {
	return c.mediaType
}

// Parameters returns detached normalized content-type parameters.
func (c ContentType) Parameters() Parameters {
	return Parameters{items: c.parameters.Items()}
}

// Equal reports whether c and other contain the same normalized key material.
func (c ContentType) Equal(other ContentType) bool {
	return c.key() == other.key()
}
