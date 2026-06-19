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

package objectreflector

import (
	"testing"
	"time"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

func waitRead(t testing.TB, ch <-chan storewatchapi.CollectionRead) storewatchapi.CollectionRead {
	t.Helper()

	select {
	case read := <-ch:
		return read
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for collection read")
		return storewatchapi.CollectionRead{}
	}
}

func waitChange(t testing.TB, ch <-chan objectstore.Change) objectstore.Change {
	t.Helper()

	select {
	case change := <-ch:
		return change
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for change")
		return objectstore.Change{}
	}
}

func requireNoChange(t testing.TB, ch <-chan objectstore.Change) {
	t.Helper()

	select {
	case change := <-ch:
		t.Fatalf("unexpected change: %#v", change)
	default:
	}
}

func TestRunPanicsOnNilContext(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	defer func() {
		if recover() == nil {
			t.Fatalf("Run(nil) did not panic")
		}
	}()
	_ = reflector.Run(nil)
}
