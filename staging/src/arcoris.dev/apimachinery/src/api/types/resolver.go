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

// Resolver resolves named structural type definitions.
//
// Resolver belongs to package types because TypeRef belongs to package types.
// Concrete catalogs belong to higher layers. This keeps structural type
// descriptors independent from catalog storage, resource registries, runtime
// schemes, codecs, converters, and global registration state.
//
// Resolver is not a recursion escape hatch. ValidateDefinition rejects
// recursive TypeDefinition graphs; recursive schemas require a future explicit
// design pass before they become part of the descriptor contract.
type Resolver interface {
	ResolveType(name TypeName) (TypeDefinition, bool)
}
