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

// NullType builds the TypeNull literal descriptor.
//
// NullType models the null literal as its own structural type. It is not a
// nullable marker for other descriptors; builders for non-null families carry
// nullability through Type flags instead.
//
// NullType deliberately has no Nullable method. TypeNull already describes the
// null literal itself and must not also carry nullable semantics.
type NullType struct {
	// header stores the descriptor kind under construction.
	header typeHeader
}

// Null returns a descriptor for the null literal type.
//
// Typical reusable declaration:
//
//	nullLiteral := Null()
func Null() NullType {
	return NullType{header: newHeader(TypeNull)}
}

// Type returns a detached Type descriptor.
func (t NullType) Type() Type {
	return typeFromHeader(t.header)
}

// typeExpr marks NullType as a sealed TypeExpr implementation.
func (t NullType) typeExpr() {}
