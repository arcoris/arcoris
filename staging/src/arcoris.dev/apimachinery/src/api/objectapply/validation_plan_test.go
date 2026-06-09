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

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

// testResolver is a minimal resolver fixture used to verify option plumbing.
type testResolver struct{}

// Resolve satisfies types.Resolver without resolving any references.
func (testResolver) Resolve(types.TypeName) (types.Definition, bool) {
	return types.Definition{}, false
}

func TestValidationPlan(t *testing.T) {
	resolver := testResolver{}
	plan := newApplier(Options{
		Resolver: resolver,
		MaxDepth: 11,
	}).validationPlan(testRequest())

	if plan.Resource.IsZero() {
		t.Fatalf("Resource is zero")
	}
	if plan.Resolver != resolver {
		t.Fatalf("Resolver was not preserved")
	}

	validator, ok := plan.DesiredValidator.(valuevalidation.SurfaceValidator)
	if !ok {
		t.Fatalf("DesiredValidator type = %T; want valuevalidation.SurfaceValidator", plan.DesiredValidator)
	}
	if validator.Options.MaxDepth != 11 {
		t.Fatalf("MaxDepth = %d; want 11", validator.Options.MaxDepth)
	}
}
