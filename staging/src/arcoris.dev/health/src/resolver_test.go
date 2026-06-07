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

package health

import "testing"

func TestCheckResolverFuncAdaptsFunction(t *testing.T) {
	t.Parallel()

	checker := mustCheck(t, "storage", Healthy("storage"))
	resolver := CheckResolverFunc(func(target Target) (CheckSet, error) {
		return NewCheckSet(target, checker)
	})

	set, err := resolver.ResolveChecks(TargetReady)
	if err != nil {
		t.Fatalf("ResolveChecks() = %v, want nil", err)
	}
	if set.Target() != TargetReady || set.Len() != 1 || !set.Has("storage") {
		t.Fatalf("set = %+v, want ready storage", set)
	}
}

func TestCheckResolverFuncNilPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("ResolveChecks on nil CheckResolverFunc did not panic")
		}
	}()

	var resolver CheckResolverFunc
	_, _ = resolver.ResolveChecks(TargetReady)
}
