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

package fieldpath

// LiteralKind identifies the scalar payload stored in a selector literal.
type LiteralKind uint8

const (
	// LiteralInvalid is the zero literal kind and is never valid selector data.
	LiteralInvalid LiteralKind = iota
	// LiteralBool stores a boolean selector literal.
	LiteralBool
	// LiteralInteger stores one exact integer from the int64 ∪ uint64 domain.
	LiteralInteger
	// LiteralString stores a string selector literal.
	LiteralString
)

// IsValid reports whether k identifies a supported selector literal kind.
func (k LiteralKind) IsValid() bool {
	return k >= LiteralBool && k <= LiteralString
}

// String returns a stable diagnostic name for k.
func (k LiteralKind) String() string {
	switch k {
	case LiteralInvalid:
		return "invalid"
	case LiteralBool:
		return "bool"
	case LiteralInteger:
		return "integer"
	case LiteralString:
		return "string"
	default:
		return "unknown"
	}
}
