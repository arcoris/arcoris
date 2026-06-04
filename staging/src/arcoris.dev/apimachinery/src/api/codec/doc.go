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

// Package codec defines common contracts for API document codecs.
//
// The package is the format-independent codec foundation for ARCORIS API
// machinery for already-configured codec implementations. It defines
// open-world codec formats, media types, document targets, typed encoder and
// decoder interfaces, optional streaming interfaces, full byte/stream
// convenience interfaces, codec metadata, and structured diagnostics.
//
// Package codec does not implement JSON, YAML, CBOR, or any other concrete wire
// format. Concrete implementations live in packages such as api/codecjson,
// api/codecyaml, and api/codeccbor. Codec registration and lookup live in a
// separate api/codecregistry package.
//
// Concrete codec packages own their typed configuration, defaults, validation,
// and effective runtime policy. Runtime callers receive configured codec
// capabilities and call Decode or Encode methods without per-call option bags.
// Different policies should be represented by different configured codec
// instances or by a higher-level selection/profile plan.
//
// Codecs transform bytes or streams to and from API document models such as
// api/value.Value, value-backed api/object envelopes, and
// api/objectownership.Document. Codecs do not validate values against api/types
// descriptors, perform object validation, apply objects, compute field
// ownership conflicts, default, prune, convert API versions, perform resource
// catalog lookup, access storage, run admission, or execute runtime/server
// lifecycle behavior.
//
// Terminology:
//
//   - BaseCodec is the metadata-only base interface shared by codec
//     implementations.
//   - ValueCodec, ObjectCodec, and ObjectOwnershipCodec are byte-slice target
//     capabilities.
//   - Codec is the full byte-slice codec interface for all v1 targets.
//   - ValueStreamCodec, ObjectStreamCodec, and ObjectOwnershipStreamCodec are
//     streaming target capabilities.
//   - StreamingCodec is the full stream codec interface for all v1 targets.
//
// Codec does not imply streaming support. StreamingCodec does not imply
// byte-slice support. Implementations may expose either surface or both.
//
// Codec diagnostic paths are syntactic document locations, not api/fieldpath
// paths. Implementations should use "$" for the encoded document root and
// locations such as "$.desired.conditions[0].status" for decoded input shape
// problems. Package metadata diagnostics use paths such as "codec.info.format".
//
// Codec implementations should preserve the difference between absent fields,
// present nulls, present empty objects, present empty arrays, and present empty
// strings whenever the format can represent those distinctions. Later
// validation, apply, and field-set extraction layers depend on that shape.
//
// Codec implementations must not silently lose numeric precision. Formats with
// arbitrary or textual numbers must preserve enough information for later value
// and type validation, and must not parse arbitrary input numbers through
// float64 when that can change the value.
//
// Formats and media types are open-world so custom codecs can exist. Targets
// are closed-world in v1 because they name API document models understood by the
// framework.
package codec
