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

package valuevalidation

import "arcoris.dev/apimachinery/api/fieldpath"

// validateIntegerLimits checks inclusive lower and upper bounds for a normalized integer.
func validateIntegerLimits[T integerValue](
	run *validator,
	path fieldpath.Path,
	got T,
	limits integerLimits[T],
) {
	if limits.lower.set && got < limits.lower.value {
		run.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonBelowMinimum,
			"integer %d is below minimum %d",
			got,
			limits.lower.value,
		)
	}

	if limits.upper.set && got > limits.upper.value {
		run.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonAboveMaximum,
			"integer %d is above maximum %d",
			got,
			limits.upper.value,
		)
	}
}
