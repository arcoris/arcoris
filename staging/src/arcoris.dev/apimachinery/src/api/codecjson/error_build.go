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

package codecjson

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// errorAt creates a direct structured JSON codec diagnostic.
func errorAt(path jsonPath, local error, codecErr error, reason ErrorReason, detail string) error {
	return &Error{
		Record:   diagnostic.NewRecord(path.String(), local, reason, detail),
		CodecErr: codecErr,
	}
}

// errorfAt creates a direct structured diagnostic with formatted detail.
func errorfAt(
	path jsonPath,
	local error,
	codecErr error,
	reason ErrorReason,
	format string,
	args ...any,
) error {
	return errorAt(path, local, codecErr, reason, fmt.Sprintf(format, args...))
}

// wrapAt creates a structured JSON codec diagnostic preserving a lower cause.
func wrapAt(path jsonPath, local error, codecErr error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Record:   diagnostic.WrapRecord(path.String(), local, reason, detail, cause),
		CodecErr: codecErr,
	}
}
