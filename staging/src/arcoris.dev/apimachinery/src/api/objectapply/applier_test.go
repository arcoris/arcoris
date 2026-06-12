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

import "testing"

func TestNewApplier(t *testing.T) {
	applier := New(Options{MaxDepth: 5, Force: true})

	if applier.opts.MaxDepth != 5 {
		t.Fatalf("MaxDepth = %d; want 5", applier.opts.MaxDepth)
	}
	if !applier.opts.Force {
		t.Fatalf("Force = false; want true")
	}
}

func TestApplierApplyMatchesPackageApply(t *testing.T) {
	req := testRequest()
	opts := Options{}

	fromFunction, err := Apply(req, opts)
	requireNoError(t, err)

	fromApplier, err := New(opts).Apply(req)
	requireNoError(t, err)

	requireStringMember(t, fromApplier.Object.Desired, "image", "api:v2")
	requireStringMember(t, fromFunction.Object.Desired, "image", "api:v2")
	requireSet(t, fromApplier.Desired.AppliedFields, "$.image")
	requireSet(t, fromFunction.Desired.AppliedFields, "$.image")
}

func TestApplierCanBeReused(t *testing.T) {
	applier := New(Options{})

	first, err := applier.Apply(testRequest())
	requireNoError(t, err)
	requireStringMember(t, first.Object.Desired, "image", "api:v2")

	req := testRequest()
	req.Applied = appliedObject(obj(member("replicas", str("4"))))
	second, err := applier.Apply(req)
	requireNoError(t, err)
	requireStringMember(t, second.Object.Desired, "replicas", "4")
}

func TestApplierStoresOptionsByValue(t *testing.T) {
	opts := Options{MaxDepth: 5}
	applier := New(opts)
	opts.MaxDepth = 99

	if applier.opts.MaxDepth != 5 {
		t.Fatalf("MaxDepth = %d; want 5", applier.opts.MaxDepth)
	}
}
