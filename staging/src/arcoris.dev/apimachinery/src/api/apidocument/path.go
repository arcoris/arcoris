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

const (
	// pathSeparator separates nested field names in diagnostic document paths.
	pathSeparator = "."
)

// Path is a stable diagnostic path into an API document shape.
//
// Path is intentionally not a JSONPath and not a fieldownership path. It names
// logical API document fields such as object.metadata.labels or
// ownership.metadata.annotations.
type Path string

// String returns the stable diagnostic path text.
func (p Path) String() string {
	return string(p)
}

// IsZero reports whether p is the zero document path.
func (p Path) IsZero() bool {
	return p == ""
}

// Child returns the nested path for field below p.
func (p Path) Child(field FieldName) Path {
	if p.IsZero() {
		return Path(field.String())
	}

	return Path(p.String() + pathSeparator + field.String())
}
