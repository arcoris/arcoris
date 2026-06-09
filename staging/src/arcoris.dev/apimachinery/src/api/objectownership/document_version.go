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

// DocumentVersion identifies the object ownership document shape.
//
// Document versions are explicit even though Document is an in-memory
// representation. Higher layers can therefore reject or deliberately migrate
// future ownership document shapes instead of guessing from field presence.
type DocumentVersion string

const (
	// DocumentVersionV1 is the first object ownership document shape.
	DocumentVersionV1 DocumentVersion = "v1"
)

// String returns the raw version text.
func (v DocumentVersion) String() string {
	return string(v)
}

// IsZero reports whether the version is absent.
func (v DocumentVersion) IsZero() bool {
	return v == ""
}

// IsSupported reports whether document conversion and validation can interpret v.
func (v DocumentVersion) IsSupported() bool {
	return v == DocumentVersionV1
}
