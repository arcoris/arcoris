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

package jitter

import (
	"testing"
	"time"
)

func TestRandomSourceFuncRejectsNilFunction(t *testing.T) {
	var source RandomSourceFunc

	mustPanicWith(t, errNilRandomSourceFunc, func() {
		_ = source.NewRandom()
	})
}

func TestRandomSourceFuncRejectsNilRandom(t *testing.T) {
	mustPanicWith(t, errNilRandom, func() {
		_ = RandomSourceFunc(func() RandomGenerator { return nil }).NewRandom()
	})
}

func TestRandomSourceFuncReturnsGenerator(t *testing.T) {
	random := RandomSourceFunc(func() RandomGenerator { return fixedRandom(7) }).NewRandom()

	if got := random.Int63(); got != 7 {
		t.Fatalf("Int63() = %d, want 7", got)
	}
}

func TestRandomFuncRejectsNilFunction(t *testing.T) {
	var random RandomFunc

	mustPanicWith(t, errNilRandomFunc, func() {
		_ = random.Int63()
	})
}

func TestRandomFuncReturnsFunctionValue(t *testing.T) {
	random := RandomFunc(func() int64 { return 11 })

	if got := random.Int63(); got != 11 {
		t.Fatalf("Int63() = %d, want 11", got)
	}
}

func TestDefaultRandomSourceReturnsGenerator(t *testing.T) {
	random := defaultRandomSource().NewRandom()
	if random == nil {
		t.Fatal("defaultRandomSource().NewRandom() returned nil")
	}
}

func TestRandomDurationInclusiveRejectsNilRandom(t *testing.T) {
	mustPanicWith(t, errNilRandom, func() {
		_ = randomDurationInclusive(nil, time.Second)
	})
}

func TestRandomDurationInclusiveUsesClosedRange(t *testing.T) {
	if got := randomDurationInclusive(fixedRandom(0), time.Second); got != 0 {
		t.Fatalf("randomDurationInclusive(0) = %s, want 0", got)
	}
	if got := randomDurationInclusive(fixedRandom(int64(time.Second)), time.Second); got != time.Second {
		t.Fatalf("randomDurationInclusive(max) = %s, want %s", got, time.Second)
	}
	if got := randomDurationInclusive(fixedRandom(99), 0); got != 0 {
		t.Fatalf("randomDurationInclusive(_, 0) = %s, want 0", got)
	}
}

func TestRandomDurationInclusiveHandlesMaxDuration(t *testing.T) {
	const maxInt63 = int64(1<<63 - 1)

	if got := randomDurationInclusive(fixedRandom(maxInt63), maxDuration); got != maxDuration {
		t.Fatalf("randomDurationInclusive(maxInt63, maxDuration) = %s, want %s", got, maxDuration)
	}
}

func TestRandomOffsetInclusiveDelegatesToDurationDraw(t *testing.T) {
	if got := randomOffsetInclusive(fixedRandom(5), 10*time.Nanosecond); got != 5*time.Nanosecond {
		t.Fatalf("randomOffsetInclusive() = %s, want 5ns", got)
	}
}
