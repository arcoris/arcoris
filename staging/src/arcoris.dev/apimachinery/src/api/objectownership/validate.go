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

import "arcoris.dev/apimachinery/api/fieldownership"

// Validate checks ownership state shape.
//
// The current public constructors make invalid State values difficult to build,
// but validation remains a package boundary for stores, codecs, tests, and any
// future internal constructors. The function is fail-fast and reports the first
// malformed surface with a stable diagnostic path.
func Validate(state State) error {
	if err := validateSurfaceState(pathStateDesired, ErrorReasonInvalidDesired, state.Desired()); err != nil {
		return err
	}
	if err := validateSurfaceState(pathStateObserved, ErrorReasonInvalidObserved, state.Observed()); err != nil {
		return err
	}
	if err := validateSurfaceState(pathStateMetadataLabels, ErrorReasonInvalidMetadataLabels, state.Metadata().Labels()); err != nil {
		return err
	}
	if err := validateSurfaceState(pathStateMetadataAnnotations, ErrorReasonInvalidMetadataAnnotations, state.Metadata().Annotations()); err != nil {
		return err
	}

	return nil
}

// validateSurfaceState attaches object-surface context to fieldownership
// structural validation.
func validateSurfaceState(path string, reason ErrorReason, state fieldownership.State) error {
	if err := state.ValidateStructure(); err != nil {
		return wrapAt(
			path,
			ErrInvalidState,
			reason,
			"object ownership surface state is invalid",
			err,
		)
	}

	return nil
}
