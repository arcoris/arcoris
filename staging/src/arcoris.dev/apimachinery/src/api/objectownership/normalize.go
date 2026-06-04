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

package objectownership

// Normalize canonicalizes doc without changing ownership semantics.
//
// It sorts owners, merges duplicate owner entries, deduplicates fields, prunes
// empty entries, and writes VersionV1. It deliberately preserves shared
// ownership and explicit parent/child path pairs.
func Normalize(doc Document) (Document, error) {
	if err := Validate(doc); err != nil {
		return Document{}, err
	}

	desired, err := fieldOwnershipStateFromDocumentSurface(pathDocumentDesired, doc.Desired)
	if err != nil {
		return Document{}, err
	}

	return Document{
		Version: VersionV1,
		Desired: fieldOwnershipStateToDocumentSurface(
			desired,
		),
	}, nil
}
