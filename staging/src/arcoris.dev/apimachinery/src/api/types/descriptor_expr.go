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

// DescriptorExpr is the sealed interface implemented by package-owned
// descriptor builders.
//
// The unexported marker prevents external packages from implementing the
// interface. That keeps Descriptor as a closed structural system instead of an
// extension point for arbitrary Go validators, reflection types, runtime object
// implementations, or transport-specific schema fragments.
type DescriptorExpr interface {
	descriptorExpr()
	Descriptor() Descriptor
}

// descriptorFromExpr converts a sealed descriptor expression into a detached
// Descriptor value.
//
// A nil expression becomes the zero Descriptor so constructors stay panic-free and
// ValidateResolved can report invalid descriptor paths consistently.
func descriptorFromExpr(expr DescriptorExpr) Descriptor {
	if expr == nil {
		return Descriptor{}
	}

	return expr.Descriptor()
}
