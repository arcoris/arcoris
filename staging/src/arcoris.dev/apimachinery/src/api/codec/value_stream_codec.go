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

	"arcoris.dev/apimachinery/api/value"
)

// ValueStreamDecoder decodes one api/value.Value document from a stream.
//
// Streaming decoders are optional capabilities for implementations that can
// process large documents without first materializing a byte slice.
type ValueStreamDecoder interface {
	DecodeValueFrom(r io.Reader) (value.Value, error)
}

// ValueStreamEncoder encodes one api/value.Value document to a stream.
//
// Implementations should write exactly one document and return any write or
// formatting failure as a structured codec error when possible.
type ValueStreamEncoder interface {
	EncodeValueTo(w io.Writer, v value.Value) error
}

// ValueStreamCodec is an optional streaming codec for api/value.Value documents.
//
// A codec may implement ValueCodec without implementing ValueStreamCodec.
type ValueStreamCodec interface {
	BaseCodec
	ValueStreamDecoder
	ValueStreamEncoder
}
