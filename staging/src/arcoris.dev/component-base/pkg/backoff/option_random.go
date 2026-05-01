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

package backoff

import "math/rand"

const (
	// errNilRandomOption is the stable diagnostic text used when a random-aware
	// schedule receives a nil RandomOption.
	//
	// Random options configure schedule construction. A nil option cannot be
	// applied meaningfully and is rejected immediately so invalid configuration
	// is reported at the backoff boundary.
	errNilRandomOption = "backoff: nil random option"

	// errNilRandomConfig is the stable diagnostic text used when package-local
	// option plumbing receives a nil randomConfig.
	//
	// Public callers cannot pass randomConfig directly. The diagnostic exists so
	// internal option wiring failures fail with a stable package-local message.
	errNilRandomConfig = "backoff: nil random config"
)

// RandomOption configures non-cryptographic randomness for random schedules and
// jitter schedules.
//
// RandomOption belongs only to mechanical delay generation. It does not add
// retry policy, context handling, sleeping, metrics, tracing, logging, error
// classification, admission control, or scheduler decisions. Invalid option
// input is a programming error and panics with a package-local diagnostic
// string.
//
// Example:
//
//	schedule := backoff.Random(
//		time.Second,
//		5*time.Second,
//		backoff.WithSeed(42),
//	)
type RandomOption func(*randomConfig)

// randomConfig is the shared configuration used by random-aware schedules.
//
// Schedules store a RandomSource so each NewSequence call can create its own
// random generator. This keeps Schedule values reusable while making mutable
// random state sequence-owned.
type randomConfig struct {
	// source creates the per-sequence random generator.
	//
	// The value is non-nil after randomOptionsOf returns. It is stored on the
	// reusable schedule; the generator returned by source.NewRandom belongs to a
	// concrete Sequence.
	source RandomSource
}

// defaultRandomConfig returns the runtime default random configuration.
//
// The default configuration uses defaultRandomSource. It is appropriate for
// best-effort runtime desynchronization and not for cryptographic randomness.
func defaultRandomConfig() randomConfig {
	return randomConfig{source: defaultRandomSource()}
}

// randomOptionsOf applies opts to the default random configuration.
//
// Nil options panic because silently ignoring them would hide invalid schedule
// construction. A nil source after option application also panics. The returned
// config always contains a non-nil RandomSource.
func randomOptionsOf(opts ...RandomOption) randomConfig {
	config := defaultRandomConfig()
	for _, option := range opts {
		requireRandomOption(option, errNilRandomOption)
		option(&config)
	}
	requireRandomConfig(&config, errNilRandomConfig)
	requireRandomSource(config.source, errNilRandomSource)
	return config
}

// WithRandomSource configures schedules to create per-sequence random generators
// from source.
//
// WithRandomSource is the preferred customization point for reusable schedules.
// The source is stored on the Schedule; each Sequence receives the random
// generator returned by source.NewRandom.
//
// Example:
//
//	source := backoff.RandomSourceFunc(func() backoff.RandomGenerator {
//		return rand.New(rand.NewSource(1))
//	})
//	schedule := backoff.FullJitter(backoff.Fixed(time.Second), backoff.WithRandomSource(source))
//
// WithRandomSource panics when source is nil.
func WithRandomSource(source RandomSource) RandomOption {
	requireRandomSource(source, errNilRandomSource)

	return func(config *randomConfig) {
		requireRandomConfig(config, errNilRandomConfig)
		config.source = source
	}
}

// WithRandom configures schedules to use random for every sequence.
//
// This adapter is mainly useful in tests that need exact boundary values from a
// small deterministic random generator. Because the same random generator value
// is returned for each sequence, callers that pass mutable implementations own
// any synchronization and sharing consequences. Production reusable schedules
// should usually prefer WithRandomSource or WithSeed.
//
// Example:
//
//	schedule := backoff.Random(0, time.Second, backoff.WithRandom(myTestRandom))
//
// WithRandom panics when random is nil.
func WithRandom(random RandomGenerator) RandomOption {
	requireRandom(random, errNilRandom)

	return WithRandomSource(RandomSourceFunc(func() RandomGenerator {
		return random
	}))
}

// WithRandomFunc adapts f into a random generator used by every sequence.
//
// The function must follow the same range contract as math/rand.Int63 and
// return values in [0, 1<<63). The adapter is intentionally small and does not
// clamp invalid values.
//
// Example:
//
//	schedule := backoff.Random(0, time.Second, backoff.WithRandomFunc(func() int64 {
//		return 0
//	}))
//
// WithRandomFunc panics when f is nil.
func WithRandomFunc(f func() int64) RandomOption {
	if f == nil {
		panic(errNilRandomFunc)
	}

	return WithRandom(RandomFunc(f))
}

// WithSeed configures schedules to create a fresh deterministic pseudo-random
// generator for each sequence.
//
// Each NewSequence call receives rand.New(rand.NewSource(seed)). Independent
// sequences from the same schedule therefore produce the same deterministic
// stream while avoiding shared mutable random state.
//
// Example:
//
//	schedule := backoff.Random(time.Second, 5*time.Second, backoff.WithSeed(42))
//	left := schedule.NewSequence()
//	right := schedule.NewSequence()
//	// left and right are deterministic and independent.
func WithSeed(seed int64) RandomOption {
	return WithRandomSource(RandomSourceFunc(func() RandomGenerator {
		return rand.New(rand.NewSource(seed))
	}))
}
