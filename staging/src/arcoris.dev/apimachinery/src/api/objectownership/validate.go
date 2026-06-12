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

// Validate checks document shape without requiring normalized ordering.
//
// Valid raw documents may contain duplicate owners, duplicate fields, unsorted
// entries, and empty field entries. Normalize owns canonicalization.
func Validate(doc Document) error {
	if err := validateDocumentVersion(doc.Version); err != nil {
		return err
	}
	if err := validateSurface(pathDocumentDesired, doc.Desired); err != nil {
		return err
	}

	return nil
}

// validateDocumentVersion rejects missing and unknown document shapes.
func validateDocumentVersion(version DocumentVersion) error {
	if version.IsZero() {
		return errorAt(
			pathDocumentVersion,
			ErrInvalidDocument,
			ErrorReasonMissingVersion,
			"document version is required",
		)
	}
	if !version.IsSupported() {
		return errorfAt(
			pathDocumentVersion,
			ErrUnsupportedVersion,
			ErrorReasonUnsupportedVersion,
			"document version %q is not supported",
			version,
		)
	}

	return nil
}

// validateSurface checks raw entries while allowing duplicate owners.
func validateSurface(path string, surface Surface) error {
	for i, entry := range surface.Entries {
		if err := validateEntry(entryPath(path, i), entry); err != nil {
			return err
		}
	}

	return nil
}

// validateEntry checks owner identity and every document path string.
func validateEntry(path string, entry Entry) error {
	if err := entry.Owner.ValidateLexical(); err != nil {
		return wrapAt(
			path+".owner",
			ErrInvalidEntry,
			ErrorReasonInvalidOwner,
			"entry owner is invalid",
			err,
		)
	}

	for i, field := range entry.Fields {
		if err := validatePath(fieldPath(path, i), field); err != nil {
			return err
		}
	}

	return nil
}

// validatePath checks one canonical document path.
func validatePath(path string, p Path) error {
	_, err := parsePath(path, p)
	return err
}
