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

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

const (
	// pathState is the root diagnostic path for whole-state invariants.
	pathState = "state"
	// pathStateDesired identifies the Desired ownership surface.
	pathStateDesired = "state.desired"
	// pathStateObserved identifies the Observed ownership surface.
	pathStateObserved = "state.observed"
	// pathStateMetadataLabels identifies metadata.labels ownership.
	pathStateMetadataLabels = "state.metadata.labels"
	// pathStateMetadataAnnotations identifies metadata.annotations ownership.
	pathStateMetadataAnnotations = "state.metadata.annotations"
)

// errorAt creates an ownership diagnostic without a lower-level cause.
func errorAt(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// errorfAt formats a detail string before creating an ownership diagnostic.
func errorfAt(path string, err error, reason ErrorReason, format string, args ...any) error {
	return errorAt(path, err, reason, fmt.Sprintf(format, args...))
}

// wrapAt creates an ownership diagnostic that preserves a lower-level cause.
func wrapAt(path string, err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Record: diagnostic.WrapRecord(path, err, reason, detail, cause),
	}
}
