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

import (
	"io"

	"arcoris.dev/apimachinery/api/objectownership"
)

// ObjectOwnershipStreamDecoder decodes object ownership documents from a stream.
//
// The stream form is for concrete codecs that can parse ownership documents
// incrementally. It still returns the in-memory objectownership.Document model.
type ObjectOwnershipStreamDecoder interface {
	DecodeObjectOwnershipFrom(r io.Reader, opts DecodeOptions) (objectownership.Document, error)
}

// ObjectOwnershipStreamEncoder encodes object ownership documents to a stream.
//
// The encoder formats the supplied document only. It does not call
// objectownership.Normalize or compute semantic ownership changes.
type ObjectOwnershipStreamEncoder interface {
	EncodeObjectOwnershipTo(w io.Writer, doc objectownership.Document, opts EncodeOptions) error
}

// ObjectOwnershipStreamCodec is an optional streaming codec for ownership docs.
//
// A codec may implement ObjectOwnershipCodec without implementing this optional
// streaming capability.
type ObjectOwnershipStreamCodec interface {
	BaseCodec
	ObjectOwnershipStreamDecoder
	ObjectOwnershipStreamEncoder
}
