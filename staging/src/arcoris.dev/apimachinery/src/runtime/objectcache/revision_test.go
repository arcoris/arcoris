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

func TestRevisionReportsReadinessSeparatelyFromRevisionValue(t *testing.T) {
	var nilCache *Cache
	_, err := nilCache.Revision()
	requireErrorIs(t, err, ErrInvalidCache)

	cache, err := New(testCollection())
	requireNoError(t, err)
	_, err = cache.Revision()
	requireErrorIs(t, err, ErrNotReady)

	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), 0)))
	revision, err := cache.Revision()
	requireNoError(t, err)
	if revision != 0 {
		t.Fatalf("ready Revision() = %s; want 0", revision)
	}
}
