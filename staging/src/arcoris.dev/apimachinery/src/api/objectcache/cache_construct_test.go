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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestNewEmptyCache(t *testing.T) {
	cache, err := New(objectstore.ListResult{})
	requireNoError(t, err)
	if cache == nil {
		t.Fatal("New() cache = nil; want cache")
	}
	if !cache.IsZero() {
		t.Fatal("IsZero() = false; want true")
	}
	if got := cache.Len(); got != 0 {
		t.Fatalf("Len() = %d; want 0", got)
	}
	if got := cache.Items(); got != nil {
		t.Fatalf("Items() = %#v; want nil", got)
	}
}
