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

// Version identifies the object ownership document shape.
//
// Versions are explicit even though Document is an in-memory representation so
// future shape changes can be handled safely by higher layers.
type Version string

const (
	// VersionV1 is the first object ownership document shape.
	VersionV1 Version = "v1"
)

// String returns the raw version text.
func (v Version) String() string {
	return string(v)
}

// IsZero reports whether the version is absent.
func (v Version) IsZero() bool {
	return v == ""
}

// IsSupported reports whether document conversion and validation can interpret v.
func (v Version) IsSupported() bool {
	return v == VersionV1
}
