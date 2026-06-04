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

import "io"

// ObjectStreamDecoder decodes one value-backed object envelope from a stream.
//
// The stream contract mirrors ObjectDecoder while avoiding an intermediate
// byte-slice requirement for implementations that support incremental parsing.
type ObjectStreamDecoder interface {
	DecodeObjectFrom(r io.Reader, opts DecodeOptions) (Object, error)
}

// ObjectStreamEncoder encodes one value-backed object envelope to a stream.
//
// The encoder writes the supplied value-backed object envelope as one document
// and must not apply object-level policy while doing so.
type ObjectStreamEncoder interface {
	EncodeObjectTo(w io.Writer, obj Object, opts EncodeOptions) error
}

// ObjectStreamCodec is an optional streaming codec for value-backed objects.
//
// A codec may implement ObjectCodec without implementing ObjectStreamCodec.
type ObjectStreamCodec interface {
	Codec
	ObjectStreamDecoder
	ObjectStreamEncoder
}
