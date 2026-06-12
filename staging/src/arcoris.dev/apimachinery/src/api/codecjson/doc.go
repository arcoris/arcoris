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

// Package codecjson implements the JSON codec for ARCORIS API documents.
//
// The package is a concrete implementation of api/codec for JSON. It supports
// value.Value, value-backed object envelopes, and objectownership.State,
// including byte-slice and stream encode/decode APIs.
//
// Object ownership JSON represents the current canonical objectownership.State
// shape directly. It does not carry domain model versions, migration markers,
// or storage-envelope metadata.
//
// Codec instances are configured at construction with jsonconfig.Config from
// api/codecjson/jsonconfig. That subpackage owns JSON-specific public
// configuration, defaults, modes, and validation. Runtime decode and encode
// methods do not accept per-call option bags.
//
// Package codecjson is descriptor-agnostic. It decodes JSON into concrete API
// value documents and object envelopes, but it does not validate values against
// api/types descriptors, default, prune, convert versions, perform resource
// catalog lookup, apply objects, compute field ownership conflicts, access
// storage, run admission, or execute runtime/server behavior.
//
// The decoder uses token-based JSON parsing into a private ordered node model.
// It preserves JSON object member order, rejects duplicate JSON object keys,
// distinguishes absent fields from explicit nulls, avoids float64-first number
// parsing, rejects invalid UTF-8 and trailing data, and reports syntactic JSON
// document paths such as "$.desired.items[0].name". These paths are not
// api/fieldpath.Path values.
//
// The generic JSON value codec supports JSON-native values. Descriptor-aware
// interpretation of strings as bytes or temporal values belongs to a future
// typed or descriptor-aware codec layer, not to this package.
//
// The package-level JSON contract is intentionally stricter than encoding/json
// defaults:
//   - diagnostic paths are syntactic JSON paths, not api/fieldpath.Path values;
//   - decoding uses an ordered private node model, not map[string]any;
//   - numbers are classified from source text, not through float64 first;
//   - object member order is preserved;
//   - duplicate object member names are rejected;
//   - invalid UTF-8 and trailing data are rejected;
//   - generic DecodeValue maps JSON strings only to value.KindString;
//   - generic EncodeValue rejects value kinds that cannot round-trip through
//     descriptor-agnostic JSON.
package codecjson
