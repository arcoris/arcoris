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
	ErrorReasonInvalidPath ErrorReason = "invalidPath"
	// ErrorReasonInvalidSyntax identifies malformed field-path text.
	ErrorReasonInvalidSyntax ErrorReason = "invalidSyntax"
	// ErrorReasonInvalidElement identifies a malformed path element.
	ErrorReasonInvalidElement ErrorReason = "invalidElement"
	// ErrorReasonInvalidSelector identifies a malformed selector.
	ErrorReasonInvalidSelector ErrorReason = "invalidSelector"
	// ErrorReasonInvalidEntry identifies a malformed selector entry.
	ErrorReasonInvalidEntry ErrorReason = "invalidEntry"
	// ErrorReasonInvalidLiteral identifies an unsupported or uninitialized literal.
	ErrorReasonInvalidLiteral ErrorReason = "invalidLiteral"
	// ErrorReasonEmptyFieldName identifies an empty field name.
	ErrorReasonEmptyFieldName ErrorReason = "emptyFieldName"
	// ErrorReasonEmptyKey identifies an empty map key.
	ErrorReasonEmptyKey ErrorReason = "emptyKey"
	// ErrorReasonNegativeIndex identifies a negative list index.
	ErrorReasonNegativeIndex ErrorReason = "negativeIndex"
	// ErrorReasonEmptySelector identifies a selector with no entries.
	ErrorReasonEmptySelector ErrorReason = "emptySelector"
	// ErrorReasonDuplicateSelectorField identifies a repeated selector field name.
	ErrorReasonDuplicateSelectorField ErrorReason = "duplicateSelectorField"
)
