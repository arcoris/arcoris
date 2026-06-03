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
	"testing"

	"arcoris.dev/apimachinery/api/valueapply"
)

func TestDesiredApplyErrorMapsConflict(t *testing.T) {
	err := desiredApplyError(valueapply.ErrConflict)

	requireErrorIs(t, err, ErrConflict)
	requireErrorIs(t, err, valueapply.ErrConflict)
}

func TestDesiredApplyErrorMapsNonConflict(t *testing.T) {
	err := desiredApplyError(valueapply.ErrInvalidValue)

	requireErrorIs(t, err, ErrDesiredApplyFailed)
	requireErrorIs(t, err, valueapply.ErrInvalidValue)
}
