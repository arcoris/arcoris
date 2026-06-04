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
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/value"
)

// Object is the value-backed API object envelope used by codec object targets.
//
// Typed object adapters belong to a future layer that can convert typed payloads
// to and from api/value.Value. The codec layer deals only with this stable
// value-backed envelope shape.
type Object = object.Object[value.Value, value.Value]

// ObjectDecoder decodes one value-backed object envelope from bytes.
//
// DecodeObject must preserve object document shape but must not run resource
// catalog lookup, object validation, version conversion, admission, or apply
// behavior.
type ObjectDecoder interface {
	DecodeObject(data []byte, opts DecodeOptions) (Object, error)
}

// ObjectEncoder encodes one value-backed object envelope to bytes.
//
// EncodeObject serializes the supplied envelope. It must not synthesize or
// update metadata such as managed fields, resource versions, or generation.
type ObjectEncoder interface {
	EncodeObject(obj Object, opts EncodeOptions) ([]byte, error)
}

// ObjectCodec is a byte-based codec for value-backed object envelopes.
//
// Implementations should also report TargetObject from Info when they expose
// this capability.
type ObjectCodec interface {
	Codec
	ObjectDecoder
	ObjectEncoder
}
