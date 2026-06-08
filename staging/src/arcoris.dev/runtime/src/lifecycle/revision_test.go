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

package lifecycle

import (
	panicassert "arcoris.dev/testutil/panic"
	"math"
	"testing"
)

func TestRevisionIsZero(t *testing.T) {
	t.Parallel()

	if !ZeroRevision.IsZero() {
		t.Fatal("ZeroRevision.IsZero() = false, want true")
	}
	if Revision(1).IsZero() {
		t.Fatal("Revision(1).IsZero() = true, want false")
	}
}

func TestRevisionNext(t *testing.T) {
	t.Parallel()

	if got := ZeroRevision.Next(); got != Revision(1) {
		t.Fatalf("ZeroRevision.Next() = %d, want 1", got)
	}
}

func TestRevisionNextPanicsOnOverflow(t *testing.T) {
	t.Parallel()

	panicassert.RequireMessage(t, errRevisionOverflow, func() {
		_ = Revision(math.MaxUint64).Next()
	})
}

func TestRevisionString(t *testing.T) {
	t.Parallel()

	if got := Revision(42).String(); got != "42" {
		t.Fatalf("Revision(42).String() = %q, want %q", got, "42")
	}
}
