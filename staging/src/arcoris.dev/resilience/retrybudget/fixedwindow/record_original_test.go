/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package fixedwindow

import (
	"testing"
	"time"
)

func TestLimiterRecordOriginal(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(0.5), WithMinRetries(0))
	prev := l.Revision()

	l.RecordOriginal()
	snap := l.Snapshot()

	if snap.Revision == prev {
		t.Fatalf("revision did not advance: %d", snap.Revision)
	}
	if snap.Value.Attempts.Original != 1 {
		t.Fatalf("Original = %d, want 1", snap.Value.Attempts.Original)
	}
	if snap.Value.Capacity.Allowed != 0 || !snap.Value.Capacity.Exhausted {
		t.Fatalf("Capacity = %+v, want allowed=0 exhausted=true", snap.Value.Capacity)
	}

	l.RecordOriginal()
	snap = l.Snapshot()
	if snap.Value.Attempts.Original != 2 {
		t.Fatalf("Original = %d, want 2", snap.Value.Attempts.Original)
	}
	if snap.Value.Capacity.Allowed != 1 || snap.Value.Capacity.Available != 1 {
		t.Fatalf("Capacity = %+v, want allowed=1 available=1", snap.Value.Capacity)
	}
}

func TestLimiterRecordOriginalRotatesWindow(t *testing.T) {
	l, clk := newTestLimiter(t, WithWindow(time.Second), WithRatio(1), WithMinRetries(0))
	l.RecordOriginal()
	prev := l.Revision()

	clk.Add(time.Second)
	l.RecordOriginal()
	snap := l.Snapshot()

	if snap.Revision == prev {
		t.Fatal("revision did not advance after rotation and record")
	}
	if snap.Value.Attempts.Original != 1 || snap.Value.Attempts.Retry != 0 {
		t.Fatalf("Attempts after rotation = %+v, want original=1 retry=0", snap.Value.Attempts)
	}
	if !snap.Value.Window.StartedAt.Equal(fixedWindowTestNow.Add(time.Second)) {
		t.Fatalf("Window.StartedAt = %s, want rotated start", snap.Value.Window.StartedAt)
	}
}
