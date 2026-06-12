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

// Package jsonconfig defines the public configuration contract for the JSON
// codec implementation.
//
// The package contains JSON-specific configuration structs, modes, defaults,
// and validation. It does not implement JSON encoding or decoding.
//
// Ownership configuration controls JSON representation of canonical
// objectownership.State values. It does not define storage schema versions or
// domain ownership-document migration policy.
//
// Configuration is resolved when a codecjson.Codec is constructed. Runtime
// encode/decode calls do not receive per-call options.
//
// The package is intentionally specific to api/codecjson. YAML, CBOR, and other
// concrete codecs should define their own configuration packages.
package jsonconfig
