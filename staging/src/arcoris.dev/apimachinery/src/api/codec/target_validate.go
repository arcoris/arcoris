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

package codec

// Validate checks that t is one of the v1 framework document targets.
//
// Unlike Format and MediaType, Target is not open-world in v1. Unknown values
// are rejected even when their text is syntactically safe.
func (t Target) Validate() error {
	return validateTargetAt(pathCodecTarget, t)
}

// validateTargetAt checks t at a caller-provided diagnostic path.
//
// Info validation uses this helper so target failures point at the relevant
// codec.info.targets[index] entry.
func validateTargetAt(path string, t Target) error {
	switch t {
	case TargetValue, TargetObject, TargetObjectOwnership:
		return nil
	case "":
		return ErrorAt(
			path,
			ErrInvalidTarget,
			ErrorReasonInvalidTarget,
			"codec target is required",
		)
	default:
		return errorfAt(
			path,
			ErrInvalidTarget,
			ErrorReasonInvalidTarget,
			"codec target %q is not supported",
			t,
		)
	}
}
