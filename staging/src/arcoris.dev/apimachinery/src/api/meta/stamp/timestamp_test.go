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

package stamp

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	if !(Timestamp{}).IsZero() {
		t.Fatal("zero Timestamp IsZero() = false")
	}

	now := NewTimestamp(time.Date(2026, 5, 30, 12, 0, 0, 123, time.UTC))
	if now.IsZero() {
		t.Fatal("non-zero Timestamp IsZero() = true")
	}
	if now.Clone().Time != now.Time {
		t.Fatalf("Clone() = %#v, want %#v", now.Clone(), now)
	}
}
