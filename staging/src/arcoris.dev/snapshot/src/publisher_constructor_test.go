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

func TestNewPublisherReturnsEmptyPublisher(t *testing.T) {
	publisher := NewPublisher[string]()

	if got := publisher.Revision(); !got.IsZero() {
		t.Fatalf("Revision() = %d, want zero", got)
	}
	if got, want := publisher.Snapshot().Value, ""; got != want {
		t.Fatalf("Snapshot().Value = %q, want %q", got, want)
	}
}

func TestNewPublisherPanicsOnNilOption(t *testing.T) {
	panicassert.RequireMessage(t, "snapshot: nil option", func() {
		_ = NewPublisher[string](nil)
	})
}
