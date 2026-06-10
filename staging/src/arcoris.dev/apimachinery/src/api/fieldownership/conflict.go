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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// Conflict records one overlap between an existing owner path and an attempted path.
//
// OwnedPath and AttemptedPath may differ because conflict detection uses
// structural overlap, not exact-only equality.
type Conflict struct {
	// Owner owns OwnedPath.
	Owner Owner

	// OwnedPath is an existing path owned by Owner.
	OwnedPath fieldpath.Path

	// AttemptedPath is the caller-supplied attempted path that overlaps OwnedPath.
	AttemptedPath fieldpath.Path
}

// ValidateStructure checks whether c is a well-formed conflict record.
//
// This method validates record structure only. It does not require OwnedPath
// and AttemptedPath to overlap; conflicts produced by State.Conflicts carry that
// stronger semantic guarantee.
func (c Conflict) ValidateStructure() error {
	return validateConflict(c)
}
