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

package admission

// NewBuiltinCatalog returns an owner-created catalog populated with built-in
// admission metadata.
//
// The returned catalog is not global and is not shared by the package. A panic
// here means admission's own built-in descriptors are internally inconsistent.
func NewBuiltinCatalog() *Catalog {
	reasons := NewBuiltinReasonRegistry()
	kinds := NewBuiltinKindRegistry()
	components := NewBuiltinComponentRegistry(kinds)

	catalog, err := NewCatalog(reasons, kinds, components)
	if err != nil {
		panic(err)
	}
	return catalog
}
