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

package objectworkqueue

import (
	"context"
	"testing"
)

func TestNormalizeContextUsesBackgroundForNil(t *testing.T) {
	if normalizeContext(nil) == nil {
		t.Fatalf("normalizeContext(nil) = nil; want context")
	}
}

func TestNormalizeContextPreservesNonNilContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey{}, "value")

	if got := normalizeContext(ctx); got != ctx {
		t.Fatalf("normalizeContext(ctx) = %#v; want original", got)
	}
}

func TestAddAndGetUseBackgroundForNilContext(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)

	requireNoError(t, queue.Add(nil, item))
	got, err := queue.Get(nil)

	requireNoError(t, err)
	requireItem(t, got, item)
}

type contextKey struct{}
