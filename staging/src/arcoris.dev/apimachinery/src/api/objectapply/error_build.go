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

package objectapply

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

const (
	// pathObject names failures that apply to the objectapply operation as a
	// whole rather than one input side.
	pathObject = "object"

	// pathObjectLive names failures rooted in the current live object.
	pathObjectLive = "object.live"

	// pathObjectApplied names failures rooted in the requested applied object.
	pathObjectApplied = "object.applied"

	// pathObjectAppliedTypeMeta names applied apiVersion/kind diagnostics.
	pathObjectAppliedTypeMeta = "object.applied.typeMeta"

	// pathObjectAppliedMetadata names applied ObjectMeta diagnostics.
	pathObjectAppliedMetadata = "object.applied.metadata"

	// pathObjectAppliedObserved names unsupported applied Observed payloads.
	pathObjectAppliedObserved = "object.applied.observed"

	// pathRequestOwner names Owner validation failures.
	pathRequestOwner = "request.owner"
	// pathRequestResource names Resource validation and version lookup failures.
	pathRequestResource = "request.resource"
	// pathObjectDesired names failures returned by value-level Desired apply.
	pathObjectDesired = "object.desired"
)

// errorAt builds an objectapply diagnostic without a nested cause.
//
// Use this for policy failures owned entirely by objectapply, such as
// unsupported metadata or observed apply.
func errorAt(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// wrapAt builds an objectapply diagnostic that preserves a lower-level cause.
//
// Use this whenever objectapply is translating objectvalidation, valueapply, or
// metadata validation errors into the objectapply taxonomy.
func wrapAt(
	path string,
	err error,
	reason ErrorReason,
	detail string,
	cause error,
) error {
	return &Error{
		Record: diagnostic.WrapRecord(path, err, reason, detail, cause),
	}
}

// errorfAt formats dynamic diagnostic details.
//
// It keeps the public error construction sites readable when diagnostic text
// needs object names, versions, or field names.
func errorfAt(
	path string,
	err error,
	reason ErrorReason,
	format string,
	args ...any,
) error {
	return errorAt(path, err, reason, fmt.Sprintf(format, args...))
}
