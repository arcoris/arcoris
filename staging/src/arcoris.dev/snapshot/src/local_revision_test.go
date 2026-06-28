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

package snapshot

import (
	"testing"

	panicassert "arcoris.dev/testutil/panic"
)

func TestRevisionIsZero(t *testing.T) {
	if !ZeroLocalRevision.IsZero() {
		t.Fatal("ZeroLocalRevision must report IsZero")
	}

	if LocalRevision(1).IsZero() {
		t.Fatal("non-zero revision reported IsZero")
	}
}

func TestRevisionNext(t *testing.T) {
	if got, want := ZeroLocalRevision.Next(), LocalRevision(1); got != want {
		t.Fatalf("Next() = %d, want %d", got, want)
	}
}

func TestRevisionNextPanicsOnOverflow(t *testing.T) {
	panicassert.RequireMessage(t, "snapshot: local revision overflow", func() {
		_ = LocalRevision(^uint64(0)).Next()
	})
}

func TestRevisionChangedSince(t *testing.T) {
	if LocalRevision(5).ChangedSince(5) {
		t.Fatal("same revision reported changed")
	}

	if !LocalRevision(5).ChangedSince(4) {
		t.Fatal("different revision reported unchanged")
	}
}
