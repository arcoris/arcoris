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

// Package apidocument defines canonical field names for ARCORIS API documents.
//
// The package is format-neutral: the same logical names are used by JSON, YAML,
// CBOR, and future codecs where the format represents API documents with named
// fields. It does not implement codecs, parsers, serializers, schema generation,
// validation, storage, resource catalog lookup, admission, or runtime behavior.
//
// These constants are the source of truth for hand-written codec
// implementations. Go struct tags in object and metadata packages remain string
// literals because struct tags cannot reference constants; tests in this package
// verify that those tags stay aligned with the constants.
package apidocument
