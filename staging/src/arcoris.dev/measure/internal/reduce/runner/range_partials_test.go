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

package runner

import (
	"strconv"
	"sync/atomic"
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestFillRangePartialsOneToOneProcessesEachRange(t *testing.T) {
	ranges := []core.Range{{Start: 0, End: 2}, {Start: 2, End: 5}, {Start: 5, End: 6}}
	partials := make([]string, len(ranges))
	fillRangePartialsOneToOne(ranges, partials, func(worker int, r core.Range, dst *string) {
		*dst = strconv.Itoa(worker) + ":" + strconv.Itoa(r.Start) + "-" + strconv.Itoa(r.End)
	})
	want := []string{"0:0-2", "1:2-5", "2:5-6"}
	for i := range want {
		if partials[i] != want[i] {
			t.Fatalf("partials = %#v, want %#v", partials, want)
		}
	}
}

func TestFillRangePartialsQueuedProcessesEveryRangeOnce(t *testing.T) {
	ranges := []core.Range{
		{Start: 0, End: 3},
		{Start: 3, End: 5},
		{Start: 5, End: 8},
		{Start: 8, End: 13},
	}
	partials := make([]int, len(ranges))
	calls := make([]atomic.Int64, len(ranges))
	fillRangePartialsQueued(ranges, partials, 2, func(_ int, r core.Range, dst *int) {
		for i, planned := range ranges {
			if planned == r {
				calls[i].Add(1)
				break
			}
		}
		*dst = r.Len()
	})
	for i, r := range ranges {
		if got := calls[i].Load(); got != 1 {
			t.Fatalf("range %d processed %d times, want 1", i, got)
		}
		if partials[i] != r.Len() {
			t.Fatalf("partials[%d] = %d, want %d", i, partials[i], r.Len())
		}
	}
}

func TestFillRangePartialsQueuedKeepsPartialsIndexedByRange(t *testing.T) {
	ranges := []core.Range{{Start: 10, End: 11}, {Start: 20, End: 22}, {Start: 30, End: 33}}
	partials := make([]int, len(ranges))
	fillRangePartialsQueued(ranges, partials, 1, func(_ int, r core.Range, dst *int) {
		*dst = r.Start
	})
	want := []int{10, 20, 30}
	for i := range want {
		if partials[i] != want[i] {
			t.Fatalf("partials = %#v, want %#v", partials, want)
		}
	}
}

func TestFillRangePartialsQueuedPublishesLocalPartialsOverDirtySlots(t *testing.T) {
	type partial struct {
		Values []int
	}
	ranges := []core.Range{{Start: 1, End: 2}, {Start: 2, End: 3}, {Start: 3, End: 4}}
	partials := []partial{
		{Values: []int{100}},
		{Values: []int{200}},
		{Values: []int{300}},
	}
	fillRangePartialsQueued(ranges, partials, 1, func(_ int, r core.Range, dst *partial) {
		dst.Values = append(dst.Values, r.Start)
	})
	for i, r := range ranges {
		if len(partials[i].Values) != 1 || partials[i].Values[0] != r.Start {
			t.Fatalf("partials[%d] = %#v, want only %d", i, partials[i].Values, r.Start)
		}
	}
}

func TestReduceBalancedRangePartialsMergesByRangeIndex(t *testing.T) {
	var scratch core.Scratch[string]
	got, ok := reduceBalancedRangePartials(
		9,
		core.Options{Workers: 3, MinItemsPerWorker: 1, Strategy: core.StrategyBalanced, MergeMode: core.MergeLinear},
		&scratch,
		func(_ int, r core.Range, dst *string) {
			*dst = strconv.Itoa(r.Start)
		}, func(dst *string, src string) {
			*dst += src
		})
	if !ok {
		t.Fatal("reduceBalancedRangePartials returned false for non-empty input")
	}
	if got != "036" {
		t.Fatalf("reduceBalancedRangePartials() = %q, want range-index order %q", got, "036")
	}
	if len(scratch.Partials) != 3 {
		t.Fatalf("partials = %d, want one partial per planned balanced range", len(scratch.Partials))
	}
}
