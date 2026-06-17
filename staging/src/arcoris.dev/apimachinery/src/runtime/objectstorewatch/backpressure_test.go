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
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestSlowWatcherLosesContinuityOnOverflow(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(1))
	stream := watchAfter(t, store, testCollection(), 0)

	createObject(t, store, testKey("system", 1), "one")
	createObject(t, store, testKey("system", 2), "two")

	_, err := stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, ErrStreamOverflow)
}
