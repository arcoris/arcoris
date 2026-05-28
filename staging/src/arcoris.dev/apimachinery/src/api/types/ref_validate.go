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

// validateRef checks TypeRef syntax, optional resolver lookup, and cycles.
func validateRef(t Type, resolver Resolver, path string, resolving map[TypeName]bool) error {
	name := t.ref.name
	if !name.IsValid() {
		return typeError(path, ErrInvalidTypeReference)
	}
	if resolver == nil {
		return nil
	}
	if resolving[name] {
		return typeError("ref("+name.String()+")", ErrInvalidTypeReference)
	}
	def, ok := resolver.ResolveType(name)
	if !ok {
		return typeError("ref("+name.String()+")", ErrUnknownTypeReference)
	}
	next := copyResolving(resolving)
	next[name] = true
	return validateType(def.Type(), resolver, "ref("+name.String()+")", next)
}
