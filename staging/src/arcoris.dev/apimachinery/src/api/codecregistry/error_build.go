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

package codecregistry

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// errorAt creates a structured registry diagnostic without a nested cause.
func errorAt(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// errorfAt creates a structured registry diagnostic with formatted detail text.
func errorfAt(path string, err error, reason ErrorReason, format string, args ...any) error {
	return errorAt(path, err, reason, fmt.Sprintf(format, args...))
}

// wrapAt creates a structured registry diagnostic preserving a nested cause.
func wrapAt(path string, err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Record: diagnostic.WrapRecord(path, err, reason, detail, cause),
	}
}
