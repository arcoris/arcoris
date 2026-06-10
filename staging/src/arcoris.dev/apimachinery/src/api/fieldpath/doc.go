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

// Package fieldpath defines structured semantic paths for ARCORIS API payload
// locations.
//
// A Path is interpreted relative to a root surface chosen by the caller. The
// root may be desired state, observed state, metadata, a value subtree, a test
// fixture payload, or another semantic surface. The fieldpath package only
// represents the path below that root; it does not assign root meaning.
//
// Field elements address fixed object or record fields. Key elements address
// dynamic map entries. Index elements address ordered list positions. Selector
// elements address associative-list entries by stable identity fields.
//
// For example:
//
//	$.spec.replicas
//	$.metadata.labels["app"]
//	$.containers[0].image
//	$.conditions[{"type":"Ready"}].status
//
// Fieldpath canonical text is not JSONPath. It is not a query language and does
// not support wildcards, predicates, filters, recursive descent, or match
// expressions. In canonical text, quoted dot-form fields remain field elements:
//
//	$."x-y"
//
// Bracket-quoted strings are dynamic map keys:
//
//	$["x-y"]
//
// FieldName is a fieldpath-level semantic field step name. It is not
// api/types.FieldName and does not own descriptor field-name grammar. MapKey is
// a fieldpath-level dynamic map-key step. It is not a label key, annotation key,
// storage key, document field name, or selector predicate.
//
// Selectors are exact associative-list identities, not filters, queries,
// predicates, or partial matches.
//
// ValidateStructure checks structural well-formedness only. It does not
// validate paths against descriptors, object schemas, ownership policy, apply
// strategy, storage layout, or wire format.
//
// CanonicalText returns the stable text form emitted by the package. String
// delegates to CanonicalText as diagnostic convenience. ParseCanonical accepts
// only the canonical grammar emitted by CanonicalText.
//
// Builder is a mutable local helper for allocation-conscious traversal. It is
// not concurrency-safe and must not be stored globally. Use Path to produce an
// immutable detached Path.
//
// Set is a canonical sorted collection of structurally valid paths. NewSet
// deduplicates exact paths only; it does not compact descendants. Use
// CompactSubtrees explicitly when subtree-marker semantics are intended.
//
// Package fieldpath does not inspect api/value payloads, read api/types
// descriptors, decode wire formats, validate values, compare payloads, apply
// changes, or manage field ownership. Descriptor-aware callers decide which
// element kind to append while interpreting values through resource schemas.
package fieldpath
