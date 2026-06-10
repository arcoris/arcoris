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

// Package identity defines concrete object metadata identity primitives.
//
// The package owns lexical metadata names, name prefixes, namespace
// markers, opaque UIDs, namespace/name object names, UID-bearing object
// identities, and typed object references. It does not define storage keys,
// route keys, cache keys, resource collection keys, watch keys, authorization
// resource strings, REST paths, UID generation policy, namespace defaulting,
// scope enforcement, admission behavior, or object lifecycle behavior.
//
// API group, version, kind, and resource identity remains in
// arcoris.dev/apimachinery/api/identity. ObjectReference composes API
// group/version/kind with namespace/name when a reference is name-only.
// ObjectIdentityReference adds UID when a reference must be pinned to one
// concrete object incarnation. Owner lifecycle policy lives in api/meta/owner,
// not in this raw identity package.
//
// Empty Namespace means namespace absence. It does not mean a default namespace,
// and this package never applies defaulting. Empty Name and empty UID are
// absence markers and are invalid where a complete object name or object
// identity is required.
//
// Metadata identity values are small comparable Go values. Trusted declarations
// and tests may use direct literals, but trust boundaries should use Parse,
// ValidateLexical, MarshalText, MarshalJSON, or CanonicalText. String returns
// diagnostic text without validation. CanonicalText validates first and returns
// an error for malformed direct literals.
//
// Scope-aware validation belongs above this package because resource scope
// comes from resource definitions. UID generation also belongs above this
// package; UID is only an opaque lexically safe token here.
package identity
