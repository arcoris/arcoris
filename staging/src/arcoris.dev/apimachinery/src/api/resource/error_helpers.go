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

package resource

import "fmt"

// definitionError creates a resource-family diagnostic.
//
// It is used for invariants owned by Definition itself, such as family identity,
// scope, duplicate versions, and version-set exposure rules.
func definitionError(path string, reason ErrorReason, detail string) error {
	return &Error{Path: path, Err: ErrInvalidDefinition, Reason: reason, Detail: detail}
}

// definitionErrorf formats a resource-family diagnostic detail.
func definitionErrorf(path string, reason ErrorReason, format string, args ...any) error {
	return definitionError(path, reason, fmt.Sprintf(format, args...))
}

// versionError creates a version-level diagnostic.
//
// It is used for invariants owned by VersionDefinition, such as version
// identity and Desired/Observed descriptor shape.
func versionError(path string, reason ErrorReason, detail string) error {
	return &Error{Path: path, Err: ErrInvalidVersion, Reason: reason, Detail: detail}
}

// scopeError creates a Scope diagnostic at the only scope scalar path.
//
// Scope parsing, validation, and scalar decoding all use the same path so
// callers do not need to special-case identical scope failures by origin.
func scopeError(reason ErrorReason, detail string) error {
	return &Error{Path: pathScope, Err: ErrInvalidScope, Reason: reason, Detail: detail}
}

// nestedDefinitionError preserves a lower-level diagnostic under a definition
// validation failure.
func nestedDefinitionError(path string, reason ErrorReason, detail string, cause error) error {
	return &Error{Path: path, Err: ErrInvalidDefinition, Reason: reason, Detail: detail, Cause: cause}
}

// nestedVersionError preserves a lower-level diagnostic under a version
// validation failure.
func nestedVersionError(path string, reason ErrorReason, detail string, cause error) error {
	return &Error{Path: path, Err: ErrInvalidVersion, Reason: reason, Detail: detail, Cause: cause}
}

// invalidJSON creates a JSON scalar decoding diagnostic.
func invalidJSON(path string, detail string, cause error) error {
	return &Error{Path: path, Err: ErrInvalidJSON, Reason: ErrorReasonInvalidJSON, Detail: detail, Cause: cause}
}

// nilReceiver creates the standard nil decoder receiver diagnostic.
func nilReceiver(path string) error {
	return &Error{
		Path:   path,
		Err:    ErrNilReceiver,
		Reason: ErrorReasonNilReceiver,
		Detail: detailDecodeTargetMustBeNonNil,
	}
}
