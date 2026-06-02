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

package types

// typeHeader is the shared working state carried by public builders.
//
// Builders keep only the TypeCode, common flags, and their exact local payload.
// They do not keep a full Type while building. Type remains the normalized
// descriptor/IR assembled at Type() boundaries.
type typeHeader struct {
	// code selects the descriptor kind being built.
	code TypeCode
	// flags records descriptor-wide builder flags such as nullability.
	flags typeFlags
}

// newHeader creates builder header state for code.
func newHeader(code TypeCode) typeHeader {
	return typeHeader{code: code}
}

// withNullable returns a copy of h that admits null values.
func (h typeHeader) withNullable() typeHeader {
	h.flags |= typeFlagNullable

	return h
}

// typeFromHeader creates the normalized descriptor shell for a builder.
func typeFromHeader(h typeHeader) Type {
	return Type{code: h.code, flags: h.flags}
}
