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

package objectwatch

import "testing"

func TestPackageContinuityContractScenario(t *testing.T) {
	validator := mustValidator(t, Start{Mode: StartAfterRevision, Revision: 10})

	requireNoError(t, validator.Accept(mustChangedEvent(t, 11)))
	requireNoError(t, validator.Accept(mustBookmarkEvent(t, 12)))
	requireNoError(t, validator.Accept(mustChangedEvent(t, 13)))
	requireNoError(t, validator.Accept(mustRestartEvent(t, RestartContinuityLost, 13)))

	err := validator.Accept(mustChangedEvent(t, 14))
	requireErrorIs(t, err, ErrClosed)
}
