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

import "testing"

func TestNewApplierStoresOptions(t *testing.T) {
	got := newApplier(Options{MaxDepth: 7, Force: true})

	if got.opts.MaxDepth != 7 {
		t.Fatalf("MaxDepth = %d; want 7", got.opts.MaxDepth)
	}
	if !got.opts.Force {
		t.Fatalf("Force = false")
	}
}

func TestNewStoresOptions(t *testing.T) {
	got := New(Options{MaxDepth: 7, Force: true})

	if got.opts.MaxDepth != 7 {
		t.Fatalf("MaxDepth = %d; want 7", got.opts.MaxDepth)
	}
	if !got.opts.Force {
		t.Fatalf("Force = false")
	}
}

func TestApplierApply(t *testing.T) {
	got, err := New(Options{}).Apply(specRequest(owner("user")))
	requireNoError(t, err)

	requireStringMember(t, got.Value, "image", "api:v2")
}
