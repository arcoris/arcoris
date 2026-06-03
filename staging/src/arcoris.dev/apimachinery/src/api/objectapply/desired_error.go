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
	"errors"

	"arcoris.dev/apimachinery/api/valueapply"
)

// desiredApplyError maps valueapply failures into objectapply classifications.
//
// Conflict errors are promoted to objectapply.ErrConflict so object-level
// callers can branch without knowing the value-level package. All other
// valueapply failures remain discoverable as causes under ErrDesiredApplyFailed.
func desiredApplyError(err error) error {
	if errors.Is(err, valueapply.ErrConflict) {
		return wrapAt(
			pathObjectDesired,
			ErrConflict,
			ErrorReasonConflict,
			"desired fields conflict with existing ownership",
			err,
		)
	}

	return wrapAt(
		pathObjectDesired,
		ErrDesiredApplyFailed,
		ErrorReasonDesiredApplyFailed,
		"desired apply failed",
		err,
	)
}
