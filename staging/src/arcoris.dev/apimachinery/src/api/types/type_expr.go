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

// TypeExpr is the sealed interface implemented by package-owned type builders.
//
// The unexported marker prevents external packages from implementing the
// interface. That keeps Type as a closed structural system instead of an
// extension point for arbitrary Go validators, reflection types, runtime object
// implementations, or transport-specific schema fragments.
type TypeExpr interface {
	typeExpr()
	Type() Type
}

// typeFromExpr converts a sealed type expression into a detached Type value.
//
// A nil expression becomes the zero Type so constructors stay panic-free and
// ValidateType can report invalid descriptor paths consistently.
func typeFromExpr(expr TypeExpr) Type {
	if expr == nil {
		return Type{}
	}

	return expr.Type()
}
