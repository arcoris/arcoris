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

import "arcoris.dev/apimachinery/api/fieldpath"

// Path is a canonical document fieldpath.Path string.
//
// The text must be the exact output of fieldpath.Path.String(). Keeping the
// document shape textual avoids exposing fieldpath.Path internals as document
// fields while still allowing lossless State reconstruction.
type Path string

// String returns the canonical path text.
func (p Path) String() string {
	return string(p)
}

// Parse reconstructs the in-memory semantic path for p.
//
// It validates that p is present, parseable, and already canonical.
func (p Path) Parse() (fieldpath.Path, error) {
	return parsePath("path", p)
}

// parsePath validates and parses canonical document path text.
//
// The diagnostic path points at the Document field that supplied p, allowing
// callers to report precise invalid document locations.
func parsePath(path string, p Path) (fieldpath.Path, error) {
	if p == "" {
		return fieldpath.Path{}, errorAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"field path is required",
		)
	}

	parsed, err := fieldpath.ParseCanonical(p.String())
	if err != nil {
		return fieldpath.Path{}, wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"field path is invalid",
			err,
		)
	}
	if parsed.String() != p.String() {
		return fieldpath.Path{}, errorfAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"field path %q is not canonical; use %q",
			p,
			parsed.String(),
		)
	}

	return parsed, nil
}
