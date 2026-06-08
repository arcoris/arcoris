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

package codecselection

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestNormalizeTargetAt(t *testing.T) {
	target, err := normalizeTargetAt("codecselection.decodeBindings[0].target", codec.TargetObject)
	requireNoError(t, err)

	if target != codec.TargetObject {
		t.Fatalf("target = %q; want %q", target, codec.TargetObject)
	}
}

func TestNormalizeTargetAtRejectsInvalidTarget(t *testing.T) {
	_, err := normalizeTargetAt("codecselection.decodeBindings[0].target", codec.Target("other"))

	requireErrorIs(t, err, ErrInvalidBinding)
	requireErrorIs(t, err, codec.ErrInvalidTarget)
	requireSelectionError(t, err, "codecselection.decodeBindings[0].target", ErrorReasonInvalidBinding)
}
