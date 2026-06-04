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

import "testing"

func TestTargetValidateAcceptsKnownTargets(t *testing.T) {
	for _, target := range []Target{TargetValue, TargetObject, TargetObjectOwnership} {
		t.Run(target.String(), func(t *testing.T) {
			requireNoError(t, target.Validate())
		})
	}
}

func TestTargetValidateRejectsZero(t *testing.T) {
	err := Target("").Validate()

	requireErrorIs(t, err, ErrInvalidTarget)
	requireCodecError(t, err, pathCodecTarget, ErrorReasonInvalidTarget)
}

func TestTargetValidateRejectsUnknown(t *testing.T) {
	err := Target("metadata").Validate()

	requireErrorIs(t, err, ErrInvalidTarget)
}
