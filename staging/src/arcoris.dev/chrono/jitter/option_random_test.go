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

import "testing"

func TestDefaultRandomConfigUsesDefaultSource(t *testing.T) {
	config := defaultRandomConfig()
	if config.source == nil {
		t.Fatal("defaultRandomConfig().source is nil")
	}
}

func TestRandomOptionsOfRejectsNilOption(t *testing.T) {
	mustPanicWith(t, errNilRandomOption, func() {
		_ = randomOptionsOf(nil)
	})
}

func TestRandomOptionsOfRejectsNilConfiguredSource(t *testing.T) {
	mustPanicWith(t, errNilRandomSource, func() {
		_ = randomOptionsOf(func(config *randomConfig) {
			config.source = nil
		})
	})
}

func TestWithRandomSourceRejectsNilSource(t *testing.T) {
	mustPanicWith(t, errNilRandomSource, func() {
		_ = WithRandomSource(nil)
	})
}

func TestWithRandomSourceStoresSource(t *testing.T) {
	source := &countingRandomSource{}
	config := randomOptionsOf(WithRandomSource(source))

	if config.source != source {
		t.Fatalf("source = %T, want original source", config.source)
	}
}

func TestWithRandomRejectsNilRandom(t *testing.T) {
	mustPanicWith(t, errNilRandom, func() {
		_ = WithRandom(nil)
	})
}

func TestWithRandomUsesSameGeneratorAdapter(t *testing.T) {
	random := fixedRandom(3)
	config := randomOptionsOf(WithRandom(random))

	if got := config.source.NewRandom(); got != random {
		t.Fatalf("NewRandom() = %v, want configured random", got)
	}
}

func TestWithRandomFuncRejectsNilFunction(t *testing.T) {
	mustPanicWith(t, errNilRandomFunc, func() {
		_ = WithRandomFunc(nil)
	})
}

func TestWithRandomFuncAdaptsFunction(t *testing.T) {
	config := randomOptionsOf(WithRandomFunc(func() int64 { return 9 }))

	if got := config.source.NewRandom().Int63(); got != 9 {
		t.Fatalf("Int63() = %d, want 9", got)
	}
}

func TestWithSeedCreatesFreshDeterministicGenerators(t *testing.T) {
	config := randomOptionsOf(WithSeed(42))

	left := config.source.NewRandom()
	right := config.source.NewRandom()

	for i := 0; i < 5; i++ {
		if got, want := left.Int63(), right.Int63(); got != want {
			t.Fatalf("draw %d mismatch: left=%d right=%d", i, got, want)
		}
	}
}
