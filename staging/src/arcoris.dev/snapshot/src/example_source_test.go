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

package snapshot_test

import (
	"fmt"
	"time"

	"arcoris.dev/snapshot"
)

func ExampleSource() {
	store := snapshot.NewStore("ready", snapshot.Identity[string])

	readValue := func(src snapshot.Source[string]) string {
		return src.Snapshot().Value
	}

	fmt.Println(readValue(store))

	// Output:
	// ready
}

func ExampleStampedSource() {
	store := snapshot.NewStore("ready", snapshot.Identity[string])

	age := func(src snapshot.StampedSource[string], now time.Time) time.Duration {
		return src.Stamped().Age(now)
	}

	updated := store.Stamped().Updated
	fmt.Println(age(store, updated.Add(2*time.Second)))

	// Output:
	// 2s
}

func ExampleRevisionSource() {
	publisher := snapshot.NewPublisher[string]()
	last := publisher.Revision()

	publisher.Publish("ready")

	changed := func(src snapshot.RevisionSource, prev snapshot.Revision) bool {
		return src.Revision().ChangedSince(prev)
	}

	fmt.Println(changed(publisher, last))

	// Output:
	// true
}
