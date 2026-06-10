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

package fieldpath

// ValidateStructure checks whether e is a structurally valid selector entry.
func (e SelectorEntry) ValidateStructure() error {
	if err := e.field.ValidateStructure(); err != nil {
		return nested(
			ErrInvalidEntry,
			ErrorReasonEmptyFieldName,
			"selector entry field name is empty",
			err,
		)
	}

	if err := e.value.ValidateStructure(); err != nil {
		return nested(
			ErrInvalidEntry,
			ErrorReasonInvalidLiteral,
			"selector entry literal is invalid",
			err,
		)
	}

	return nil
}
