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

package admission

import (
	"reflect"
	"testing"
)

// requirePanicValue verifies fn panics with want.
func requirePanicValue(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		if got != want {
			t.Fatalf("panic = %v, want %v", got, want)
		}
	}()

	fn()
}

// requireNoMethod fails if typ exposes an exported method that should not be
// part of the public admission contract.
func requireNoMethod(t *testing.T, typ reflect.Type, name string) {
	t.Helper()

	if _, ok := typ.MethodByName(name); ok {
		t.Fatalf("%s exposes %s, want removed", typ, name)
	}
}

// requireDecision verifies a constructor returned the expected valid decision.
func requireDecision(t *testing.T, got Decision, want Decision) {
	t.Helper()

	if got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}
	if !got.IsValid() {
		t.Fatalf("decision is invalid: %+v", got)
	}
}

// requireResultShape verifies a constructor returned the expected valid result
// shape without asserting domain metadata or grant values.
func requireResultShape[G any, M any](
	t *testing.T,
	result Result[G, M],
	decision Decision,
	hasGrant bool,
	hasMetadata bool,
) {
	t.Helper()

	if !result.IsValid() {
		t.Fatalf("result is invalid: %+v", result.Decision())
	}
	if got := result.Decision(); got != decision {
		t.Fatalf("Decision() = %+v, want %+v", got, decision)
	}
	if got := result.HasGrant(); got != hasGrant {
		t.Fatalf("HasGrant() = %t, want %t", got, hasGrant)
	}
	if got := result.HasMetadata(); got != hasMetadata {
		t.Fatalf("HasMetadata() = %t, want %t", got, hasMetadata)
	}
}

// matrixName formats a compact decision/result matrix case name.
func matrixName(outcome Outcome, effect Effect, grantPresent bool, metadataPresent bool) string {
	name := outcome.String() + "_" + effect.String()
	if grantPresent {
		name += "_grant"
	} else {
		name += "_no_grant"
	}
	if metadataPresent {
		name += "_metadata"
	} else {
		name += "_no_metadata"
	}
	return name
}
