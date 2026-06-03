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

package typeref

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestNewStoresResolverAndDepth(t *testing.T) {
	resolver := resolverFunc(stringResolver)
	references := New(resolver, 7)

	if references.resolver == nil {
		t.Fatalf("resolver was not stored")
	}
	if references.maxDepth != 7 {
		t.Fatalf("maxDepth = %d; want 7", references.maxDepth)
	}
	if references.active == nil {
		t.Fatalf("active map was not initialized")
	}
}

func TestEnterMarksAndClearsActiveReference(t *testing.T) {
	references := New(resolverFunc(stringResolver), 64)

	leave := references.Enter("example.Name")
	if !references.active["example.Name"] {
		t.Fatalf("reference was not marked active")
	}

	leave()
	if references.active["example.Name"] {
		t.Fatalf("reference was not cleared")
	}
}

func TestEnterCleanupIsIdempotent(t *testing.T) {
	references := New(resolverFunc(stringResolver), 64)

	leave := references.Enter(types.TypeName("example.Name"))
	leave()
	leave()

	if references.active["example.Name"] {
		t.Fatalf("reference remained active after cleanup")
	}
}
