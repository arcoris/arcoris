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

// validateNewEntry checks the entry before construction returns it.
func validateNewEntry(entry Entry) error {
	return entry.ValidateStructure()
}

// ValidateStructure checks whether e is a structurally valid ownership entry.
//
// It does not check conflicts with other owners or object-level ownership
// policy. Empty field sets are valid Entry values; State normalization prunes
// them.
func (e Entry) ValidateStructure() error {
	if err := e.owner.ValidateLexical(); err != nil {
		return wrapAt(
			"entry.owner",
			ErrInvalidEntry,
			ErrorReasonInvalidEntryOwner,
			"entry owner is invalid",
			err,
		)
	}
	if err := validateFieldsAt("entry.fields", e.fields, ErrorReasonInvalidEntryFields, "entry field path is invalid"); err != nil {
		return err
	}

	return nil
}
