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

// FieldExpr is the sealed interface implemented by package-owned field builders.
//
// The unexported marker prevents external packages from injecting arbitrary
// field implementations into object descriptors. Field builders also do not
// implement TypeExpr: a field has a name, presence, description, and type; it
// is not itself a reusable unnamed type.
type FieldExpr interface {
	fieldExpr()
	Field() FieldDescriptor
}

// fieldFromExpr converts a sealed field expression into a detached descriptor.
//
// A nil expression becomes the zero FieldDescriptor so object constructors can
// stay panic-free and descriptor validation can report the invalid field path.
func fieldFromExpr(expr FieldExpr) FieldDescriptor {
	if expr == nil {
		return FieldDescriptor{}
	}
	return expr.Field()
}
