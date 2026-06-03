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

package objectownership

// Document is the versioned in-memory object ownership document.
//
// A Document is not a codec type and does not define JSON, YAML, binary, or
// storage behavior. It stores stable semantic ownership data that a future
// higher-level codec or storage layer can serialize if needed.
type Document struct {
	// Version declares the document shape. It must be explicit so future format
	// changes can be rejected or handled deliberately.
	Version Version

	// Desired stores ownership for the object's Desired surface. v1 does not
	// model Observed or metadata ownership.
	Desired Surface
}

// IsEmpty reports whether the document contains no owned object fields.
//
// IsEmpty is not a validity check. It ignores Version and every other document
// validation concern, and only answers whether normalized ownership would be
// empty.
//
// Empty raw entries do not make a document semantically non-empty because
// Normalize prunes them before constructing State.
func (d Document) IsEmpty() bool {
	return d.Desired.IsEmpty()
}
