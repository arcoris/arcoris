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

// BoolType builds boolean descriptors.
//
// BoolType records the structural contract for boolean API values. It has no
// value constraints in this design pass and exists as a closed descriptor
// builder rather than a Go bool wrapper.
type BoolType struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
}

// Bool returns a descriptor builder for boolean values.
//
// Typical reusable declaration:
//
//	enabledType := Bool().Nullable()
func Bool() BoolType {
	return BoolType{header: newHeader(TypeBool)}
}

// Nullable returns a boolean descriptor that admits null values.
func (t BoolType) Nullable() BoolType {
	t.header = t.header.withNullable()
	return t
}

// Type returns a detached Type descriptor.
func (t BoolType) Type() Type {
	return typeFromHeader(t.header)
}

// typeExpr marks BoolType as a sealed TypeExpr implementation.
func (t BoolType) typeExpr() {}
