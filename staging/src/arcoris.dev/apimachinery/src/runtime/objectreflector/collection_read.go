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

package objectreflector

import storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"

// validateCollectionRead verifies the source returned a usable boundary for the
// exact collection this Reflector owns.
//
// CollectionRead.Validate checks the objectstorewatch value contract. The
// collection equality check is a source contract check: a source that lists a
// different collection may still return a valid CollectionRead, but consuming it
// would make the subsequent watch boundary unrelated to this Reflector's sink.
func (r *Reflector) validateCollectionRead(read storewatchapi.CollectionRead) error {
	if err := read.Validate(); err != nil {
		return err
	}
	if read.Collection() != r.collection {
		return sourceContractError("collection read belongs to %v, want %v", read.Collection(), r.collection)
	}

	return nil
}
