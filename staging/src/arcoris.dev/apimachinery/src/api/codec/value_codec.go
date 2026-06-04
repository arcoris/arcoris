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

import "arcoris.dev/apimachinery/api/value"

// ValueDecoder decodes one api/value.Value document from bytes.
//
// DecodeValue performs format decoding only. It must not validate the decoded
// value against api/types descriptors, apply defaults, prune fields, or infer
// resource-specific semantics.
type ValueDecoder interface {
	DecodeValue(data []byte, opts DecodeOptions) (value.Value, error)
}

// ValueEncoder encodes one api/value.Value document to bytes.
//
// EncodeValue serializes the concrete value model as supplied. It must not
// normalize semantic values, update field ownership, or perform descriptor
// validation as part of byte formatting.
type ValueEncoder interface {
	EncodeValue(v value.Value, opts EncodeOptions) ([]byte, error)
}

// ValueCodec is a byte-based codec for api/value.Value documents.
//
// Implementations should also report TargetValue from Info when they expose
// this capability.
type ValueCodec interface {
	BaseCodec
	ValueDecoder
	ValueEncoder
}
