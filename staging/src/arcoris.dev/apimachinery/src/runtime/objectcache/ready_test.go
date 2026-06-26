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

package objectcache

import (
	"context"
	"testing"
)

func TestReadyDistinguishesUninitializedFromReadyEmpty(t *testing.T) {
	var nilCache *Cache
	if nilCache.Ready() {
		t.Fatalf("nil Ready() = true; want false")
	}

	cache, err := New(testCollection())
	requireNoError(t, err)
	if cache.Ready() {
		t.Fatalf("new Ready() = true; want false")
	}

	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), 0)))
	if !cache.Ready() {
		t.Fatalf("ready empty cache reports false")
	}
}
