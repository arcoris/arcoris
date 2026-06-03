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

package builtin

import "arcoris.dev/admissioncatalog"

// NewCatalog returns an owner-created catalog populated with standard admission
// metadata.
//
// The returned catalog is not global and is not shared by the package. A panic
// here means this package's standard descriptors are internally inconsistent.
func NewCatalog() *admissioncatalog.Catalog {
	reasons := NewReasonRegistry()
	kinds := NewKindRegistry()
	components := NewComponentRegistry(kinds)

	catalog, err := admissioncatalog.NewCatalog(reasons, kinds, components)
	if err != nil {
		panic(err)
	}
	return catalog
}
