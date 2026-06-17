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

package objectstorewatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestWatchOptionsDefaultDisallowsProgress(t *testing.T) {
	request, err := mustBoundary(t, 1).WatchRequest(WatchOptions{})
	requireNoError(t, err)

	if request.AllowProgress {
		t.Fatalf("AllowProgress = true; want false")
	}
}

func TestWatchOptionsAllowProgress(t *testing.T) {
	request, err := mustBoundary(t, 1).WatchRequest(WatchOptions{AllowProgress: true})
	requireNoError(t, err)

	if !request.AllowProgress {
		t.Fatalf("AllowProgress = false; want true")
	}
}

func mustBoundary(t testing.TB, revision uint64) Boundary {
	t.Helper()

	boundary, err := NewBoundary(testCollection(), objectstore.Revision(revision))
	requireNoError(t, err)

	return boundary
}
