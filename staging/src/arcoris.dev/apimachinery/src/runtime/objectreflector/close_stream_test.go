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

package objectreflector

import (
	"errors"
	"testing"
)

func TestCloseStreamAndJoinPreservesExistingError(t *testing.T) {
	existing := errors.New("cycle failed")
	closeErr := errors.New("close failed")
	stream := terminalStream(existing)
	stream.closeErr = closeErr
	err := existing

	closeStreamAndJoin(&err, stream)

	requireErrorIs(t, err, existing)
	requireErrorIs(t, err, closeErr)
}

func TestCloseStreamAndJoinUsesCloseErrorWhenCycleSucceeded(t *testing.T) {
	closeErr := errors.New("close failed")
	stream := terminalStream(closeErr)
	stream.closeErr = closeErr
	var err error

	closeStreamAndJoin(&err, stream)

	requireErrorIs(t, err, closeErr)
}
