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

// validateObservedPolicy rejects applied observed data before object validation.
//
// Even resources that define Observed cannot accept Observed through apply v1.
// Observed is a live/read surface and must be updated by future status/runtime
// layers, not by Desired apply.
func validateObservedPolicy(applied ValueObject) error {
	if applied.Observed == nil {
		return nil
	}

	return errorAt(
		pathObjectAppliedObserved,
		ErrUnsupportedObservedApply,
		ErrorReasonUnsupportedObservedApply,
		"applied observed surface is unsupported",
	)
}
