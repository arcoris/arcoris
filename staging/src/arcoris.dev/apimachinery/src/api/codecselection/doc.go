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

// Package codecselection builds immutable explicit codec selection plans.
//
// The package sits above api/codec and api/codecregistry. api/codec defines
// format-independent capability contracts. api/codecregistry catalogs
// already-configured codec candidates by stable EntryID and groups them by
// format and media type. Package codecselection binds normalized request keys to
// exact EntryID values and returns typed codec capabilities for those bindings.
//
// EntryID is the selected identity. Media types and formats are not identities:
// several configured codec candidates can intentionally expose the same media
// type and format while using different limits, output ordering, strictness, or
// future profile policies. Selection therefore never chooses the first registry
// entry for a media type. Ambiguity is resolved only by construction-time
// bindings that name a concrete EntryID.
//
// ContentType and Parameters in this package are normalized key material. They
// are not raw MIME headers and they do not interpret parameter semantics such as
// canonical profiles or runtime roles. Future MIME/profile packages can parse
// protocol inputs and construct selection configs explicitly.
//
// Selector construction validates all bindings against the supplied registry:
// referenced entries must exist, declare the bound media type and target, and
// implement the requested byte or stream capability. Selectors are immutable and
// safe for concurrent lookup after construction. They do not configure codecs
// and they do not make the selected codec implementations themselves
// concurrency-safe.
//
// The package does not parse Content-Descriptor or Accept headers, implement wildcard
// matching, interpret q-values, define public/storage/import/debug roles,
// create codecs, configure codecs, validate API values against descriptors,
// apply objects, access storage, run admission, or install global defaults.
package codecselection
