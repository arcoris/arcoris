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

import "arcoris.dev/apimachinery/api/objectownership"

// ObjectOwnershipDecoder decodes object ownership documents from bytes.
//
// DecodeObjectOwnership converts a concrete format into the stable in-memory
// objectownership.Document representation. Semantic document validation remains
// the responsibility of api/objectownership, not the codec contract itself.
type ObjectOwnershipDecoder interface {
	DecodeObjectOwnership(data []byte, opts DecodeOptions) (objectownership.Document, error)
}

// ObjectOwnershipEncoder encodes object ownership documents to bytes.
//
// EncodeObjectOwnership formats an existing ownership document. It must not
// compute ownership, merge ownership, or mutate object metadata.
type ObjectOwnershipEncoder interface {
	EncodeObjectOwnership(doc objectownership.Document, opts EncodeOptions) ([]byte, error)
}

// ObjectOwnershipCodec is a byte-based codec for object ownership documents.
//
// Implementations should also report TargetObjectOwnership from Info when they
// expose this capability.
type ObjectOwnershipCodec interface {
	Codec
	ObjectOwnershipDecoder
	ObjectOwnershipEncoder
}
