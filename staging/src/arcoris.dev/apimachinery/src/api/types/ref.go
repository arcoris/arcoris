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

// RefType builds references to named structural TypeDefinition values.
//
// TypeRef is the descriptor reuse mechanism. It never represents arbitrary Go
// implementations, package globals, reflection types, or runtime object types.
// Recursive TypeDefinition graphs are not supported; recursive schemas need a
// future explicit design pass before TypeRef can carry those semantics.
type RefType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact reference target under construction.
	payload refPayload
}

// Ref returns a reference descriptor builder for name.
//
// The name parameter accepts string-like values so descriptor declarations can
// use string literals while the stored payload remains a TypeName.
//
// Typical reusable declaration:
//
//	nameRef := Ref("arcoris.meta.Name")
//	nameRef = nameRef.Nullable()
func Ref[N ~string](name N) RefType {
	return RefType{
		header:  newHeader(TypeRef),
		payload: refPayload{name: TypeName(name)},
	}
}

// Nullable returns a reference descriptor that admits null values.
func (t RefType) Nullable() RefType { t.header = t.header.withNullable(); return t }

// Type returns a detached Type descriptor.
func (t RefType) Type() Type {
	out := typeFromHeader(t.header)
	out.ref = cloneRefPayload(t.payload)
	return out
}

// typeExpr marks RefType as a sealed TypeExpr implementation.
func (t RefType) typeExpr() {}
