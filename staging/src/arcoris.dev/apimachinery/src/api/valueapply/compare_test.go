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

package valueapply

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestCompare(t *testing.T) {
	req := specRequest(owner("user"))

	got, err := newApplier(Options{}).compare(req)
	requireNoError(t, err)

	requireSet(t, got.Changed(), "$.image", "$.replicas")
}

func TestCompareOptions(t *testing.T) {
	got := newApplier(Options{MaxDepth: 11}).compareOptions()

	if got.MaxDepth != 11 {
		t.Fatalf("MaxDepth = %d; want 11", got.MaxDepth)
	}
}

func TestApplyValueCompareErrorWrapped(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       str("old"),
		Applied:    str("new"),
		Descriptor: types.Bool().Type(),
	}

	_, err := newApplier(Options{}).compare(req)

	requireErrorIs(t, err, ErrCompareFailed)
}
