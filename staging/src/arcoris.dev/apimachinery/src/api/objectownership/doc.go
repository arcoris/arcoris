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

// Package objectownership models object-level field ownership state.
//
// The package is the object-level counterpart to api/fieldownership. It stores
// ownership for object surfaces such as Desired and provides a stable,
// deterministic, versioned in-memory document representation. Higher layers may
// decide how that Document is serialized or stored.
//
// Document version 1 owns only the Desired surface. The State shape is
// intentionally private so Observed or metadata ownership can be added later
// without exposing public struct fields.
//
// State is the operational object ownership model. Document is the stable
// in-memory document representation. ToDocument converts State to Document.
// StateFromDocument converts Document to State. Normalize canonicalizes
// Document. Validate checks Document shape without performing runtime, storage,
// codec, JSON, YAML, binary, or wire-format behavior.
//
// The package does not apply objects, merge values, validate resources, inspect
// api/value payloads, serialize whole objects, mutate ObjectMeta, run
// admission, authorize request subjects, access storage, perform resource
// catalog lookup, or execute runtime lifecycle behavior. Runtime and codec
// layers are responsible for those concerns.
package objectownership
