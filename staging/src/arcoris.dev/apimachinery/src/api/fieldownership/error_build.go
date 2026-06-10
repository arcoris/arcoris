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

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// errorAt builds a structured diagnostic at a field ownership location.
func errorAt(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// errorfAt builds a structured diagnostic with formatted detail text.
func errorfAt(path string, err error, reason ErrorReason, format string, args ...any) error {
	return errorAt(
		path,
		err,
		reason,
		fmt.Sprintf(format, args...),
	)
}

// wrapAt builds a structured diagnostic that preserves a lower-level cause.
func wrapAt(path string, err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Record: diagnostic.WrapRecord(path, err, reason, detail, cause),
	}
}

// wrapPathError maps a fieldpath validation failure into this package's errors.
func wrapPathError(path fieldpath.Path, detail string, cause error) error {
	return wrapAt(
		path.CanonicalText(),
		ErrInvalidPath,
		ErrorReasonInvalidPath,
		detail,
		cause,
	)
}
