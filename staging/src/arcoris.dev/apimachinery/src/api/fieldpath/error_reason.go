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

// ErrorReason identifies one precise field-path failure.
type ErrorReason string

const (
	// ErrorReasonInvalidPath identifies a path that contains an invalid element.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"
	// ErrorReasonInvalidSyntax identifies malformed field-path text.
	ErrorReasonInvalidSyntax ErrorReason = "invalid_syntax"
	// ErrorReasonInvalidElement identifies a malformed path element.
	ErrorReasonInvalidElement ErrorReason = "invalid_element"
	// ErrorReasonInvalidSelector identifies a malformed selector.
	ErrorReasonInvalidSelector ErrorReason = "invalid_selector"
	// ErrorReasonInvalidEntry identifies a malformed selector entry.
	ErrorReasonInvalidEntry ErrorReason = "invalid_entry"
	// ErrorReasonInvalidLiteral identifies an unsupported or uninitialized literal.
	ErrorReasonInvalidLiteral ErrorReason = "invalid_literal"
	// ErrorReasonEmptyFieldName identifies an empty field name.
	ErrorReasonEmptyFieldName ErrorReason = "empty_field_name"
	// ErrorReasonEmptyMapKey identifies an empty map key.
	ErrorReasonEmptyMapKey ErrorReason = "empty_map_key"
	// ErrorReasonNegativeIndex identifies a negative list index.
	ErrorReasonNegativeIndex ErrorReason = "negative_index"
	// ErrorReasonEmptySelector identifies a selector with no entries.
	ErrorReasonEmptySelector ErrorReason = "empty_selector"
	// ErrorReasonDuplicateSelectorField identifies a repeated selector field name.
	ErrorReasonDuplicateSelectorField ErrorReason = "duplicate_selector_field"
	// ErrorReasonNonCanonicalText identifies valid path text that is not canonical.
	ErrorReasonNonCanonicalText ErrorReason = "non_canonical_text"
)
