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

package apidocument

// FieldName is a canonical field name in an ARCORIS API document.
//
// FieldName values describe logical document fields shared by JSON, YAML, CBOR,
// and future codecs. They are not JSON-specific and do not define codec
// behavior.
type FieldName string

// String returns the canonical field name text.
func (n FieldName) String() string {
	return string(n)
}

// IsZero reports whether n is the zero field name.
func (n FieldName) IsZero() bool {
	return n == ""
}
