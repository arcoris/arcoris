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

package atomicx

import (
	"sync"
	"testing"
)

type pointerTestValue struct {
	value int
}

func TestPaddedPointerZeroValueLoadsNil(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	if got := p.Load(); got != nil {
		t.Fatalf("Load() = %p, want nil", got)
	}
}

func TestPaddedPointerStoreLoad(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	value := &pointerTestValue{value: 1}

	p.Store(value)

	if got := p.Load(); got != value {
		t.Fatalf("Load() = %p, want %p", got, value)
	}
}

func TestPaddedPointerSwap(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	first := &pointerTestValue{value: 1}
	second := &pointerTestValue{value: 2}
	p.Store(first)

	if got := p.Swap(second); got != first {
		t.Fatalf("Swap() = %p, want %p", got, first)
	}
	if got := p.Load(); got != second {
		t.Fatalf("Load() = %p, want %p", got, second)
	}
}

func TestPaddedPointerCompareAndSwapSuccess(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	first := &pointerTestValue{value: 1}
	second := &pointerTestValue{value: 2}
	p.Store(first)

	if !p.CompareAndSwap(first, second) {
		t.Fatal("CompareAndSwap() = false, want true")
	}
	if got := p.Load(); got != second {
		t.Fatalf("Load() = %p, want %p", got, second)
	}
}

func TestPaddedPointerCompareAndSwapFailure(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	first := &pointerTestValue{value: 1}
	second := &pointerTestValue{value: 2}
	third := &pointerTestValue{value: 3}
	p.Store(first)

	if p.CompareAndSwap(second, third) {
		t.Fatal("CompareAndSwap() = true, want false")
	}
	if got := p.Load(); got != first {
		t.Fatalf("Load() = %p, want %p", got, first)
	}
}

func TestPaddedPointerNilTransitions(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	value := &pointerTestValue{value: 1}

	if !p.CompareAndSwap(nil, value) {
		t.Fatal("CompareAndSwap(nil, value) = false, want true")
	}
	if got := p.Load(); got != value {
		t.Fatalf("Load() = %p, want %p", got, value)
	}

	p.Store(nil)
	if got := p.Load(); got != nil {
		t.Fatalf("Load() after Store(nil) = %p, want nil", got)
	}

	p.Store(value)
	if got := p.Swap(nil); got != value {
		t.Fatalf("Swap(nil) = %p, want %p", got, value)
	}
	if got := p.Load(); got != nil {
		t.Fatalf("Load() after Swap(nil) = %p, want nil", got)
	}
}

func TestPaddedPointerConcurrentLoadStore(t *testing.T) {
	t.Parallel()

	var p PaddedPointer[pointerTestValue]
	values := []pointerTestValue{
		{value: 1},
		{value: 2},
		{value: 3},
	}

	var wg sync.WaitGroup
	for i := range values {
		wg.Add(1)
		go func(value *pointerTestValue) {
			defer wg.Done()
			for range 1000 {
				p.Store(value)
				_ = p.Load()
			}
		}(&values[i])
	}

	wg.Wait()
}
