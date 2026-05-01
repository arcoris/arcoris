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

import (
	"math/rand"
	"time"
)

const (
	// errNilRandomSource is the stable diagnostic text used when a random source
	// adapter receives a nil source.
	//
	// A random source is responsible for creating per-sequence random generators.
	// A nil source cannot provide deterministic or runtime random values and
	// indicates invalid package use rather than schedule exhaustion.
	errNilRandomSource = "backoff: nil random source"

	// errNilRandomSourceFunc is the stable diagnostic text used when a
	// RandomSourceFunc method is called on a nil function value.
	//
	// RandomSourceFunc adapts a function into RandomSource. A nil function cannot
	// create a random generator and is rejected immediately at the random-source
	// boundary.
	errNilRandomSourceFunc = "backoff: nil RandomSourceFunc"

	// errNilRandom is the stable diagnostic text used when a random source returns
	// a nil random generator.
	//
	// RandomSource.NewRandom must return a usable generator. Returning nil
	// violates the RandomSource contract and would otherwise move a configuration
	// error into a later backoff, jitter, retry, polling, reconnect, or cooldown
	// path.
	errNilRandom = "backoff: nil random"

	// errNilRandomFunc is the stable diagnostic text used when a RandomFunc method
	// is called on a nil function value.
	//
	// RandomFunc adapts a function into a random generator. A nil function cannot
	// produce pseudo-random values and is rejected immediately at the random
	// boundary.
	errNilRandomFunc = "backoff: nil RandomFunc"
)

// RandomGenerator provides non-cryptographic pseudo-random int63 values.
//
// RandomGenerator is intentionally small. The package needs only one primitive
// to build inclusive duration draws. Implementations must return values in the
// same range as math/rand.Int63:
//
//	[0, 1<<63)
//
// The interface is named RandomGenerator instead of Random because the package
// already exposes the Random constructor function for random delay schedules and
// Go package-level declarations share one namespace. Keeping the constructor
// name preserves the schedule API while still making the generator contract
// explicit.
//
// RandomGenerator values are used only for runtime desynchronization and
// deterministic tests. They MUST NOT be used for security decisions, secrets,
// cryptographic nonces, authentication tokens, randomized access control, or any
// other security-sensitive purpose.
//
// RandomGenerator implementations may be stateful. Unless a concrete
// implementation documents stronger guarantees, RandomGenerator values should be
// treated as single-owner values. Schedules should create a fresh generator for
// each Sequence when they need mutable random state.
type RandomGenerator interface {
	// Int63 returns a non-negative pseudo-random int64 value in [0, 1<<63).
	Int63() int64
}

// RandomSource creates random generators for backoff sequences.
//
// RandomSource separates reusable schedule configuration from per-sequence
// random state. A Schedule can store a RandomSource, and each NewSequence call
// can request a fresh generator from it. This prevents unrelated retry, polling,
// reconnect, or cooldown loops from accidentally sharing mutable random state.
//
// RandomSource implementations should be safe for concurrent calls to NewRandom
// unless they explicitly document weaker ownership rules. RandomGenerator values
// returned by NewRandom are single-owner by default.
type RandomSource interface {
	// NewRandom creates a random generator for one sequence owner.
	//
	// The returned generator must be non-nil.
	NewRandom() RandomGenerator
}

// RandomSourceFunc adapts a function into a RandomSource.
//
// RandomSourceFunc is useful for tests and custom adapters that need to provide
// deterministic per-sequence random generators without declaring a named type.
// The wrapped function is called every time NewRandom is invoked.
//
// A nil RandomSourceFunc is a programming error. NewRandom panics immediately
// with errNilRandomSourceFunc. A function that returns nil violates the
// RandomSource contract and panics with errNilRandom.
type RandomSourceFunc func() RandomGenerator

// NewRandom calls f and returns the random generator produced by it.
//
// NewRandom panics when f is nil or when f returns nil. Both cases are
// programming errors, not schedule exhaustion.
func (f RandomSourceFunc) NewRandom() RandomGenerator {
	if f == nil {
		panic(errNilRandomSourceFunc)
	}

	r := f()
	requireRandom(r, errNilRandom)

	return r
}

// RandomFunc adapts a function into a RandomGenerator.
//
// RandomFunc is useful for tests that need deterministic boundary behavior. The
// function must return values in [0, 1<<63), matching math/rand.Int63. This
// adapter does not clamp or repair invalid values because doing so would hide a
// broken test generator or custom random implementation.
//
// A nil RandomFunc is a programming error. Int63 panics immediately with
// errNilRandomFunc.
type RandomFunc func() int64

// Int63 calls f and returns the pseudo-random value produced by it.
//
// Int63 panics when f is nil.
func (f RandomFunc) Int63() int64 {
	if f == nil {
		panic(errNilRandomFunc)
	}

	return f()
}

// defaultRandomSource returns the package default random source.
//
// The default source delegates to the standard library package-level
// pseudo-random generator. It is suitable for runtime desynchronization and load
// spreading. It is not cryptographic randomness.
//
// The returned source is stateless. Each NewRandom call returns a runtimeRandom
// adapter over package-level math/rand functions.
func defaultRandomSource() RandomSource {
	return runtimeRandomSource{}
}

// runtimeRandomSource creates runtimeRandom adapters.
//
// runtimeRandomSource is intentionally empty. It carries no mutable state and is
// safe to copy.
type runtimeRandomSource struct{}

// NewRandom returns a RandomGenerator backed by package-level math/rand
// functions.
//
// Package-level math/rand functions are appropriate for best-effort runtime
// desynchronization. Callers that need deterministic tests should use
// RandomSourceFunc, WithSeed, or another package-local test source instead.
func (runtimeRandomSource) NewRandom() RandomGenerator {
	return runtimeRandom{}
}

// runtimeRandom adapts package-level math/rand functions to RandomGenerator.
//
// The adapter does not own a *rand.Rand instance and therefore does not expose
// mutable random state through the value itself.
type runtimeRandom struct{}

// Int63 returns a non-negative pseudo-random value from package-level math/rand.
func (runtimeRandom) Int63() int64 {
	return rand.Int63()
}

// randomDurationInclusive returns a pseudo-random duration in [0, max].
//
// The helper maps the non-negative Int63 value into the requested closed range
// with modulo arithmetic. That keeps the helper total even for deterministic
// test generators that repeatedly return the same valid Int63 value. The result
// is intentionally non-cryptographic and may have a small modulo bias for ranges
// that do not divide 1<<63 evenly, which is acceptable for backoff
// desynchronization and deterministic package tests.
//
// The function assumes max is non-negative. Public constructors and package
// helpers should validate that invariant before calling it.
func randomDurationInclusive(r RandomGenerator, max time.Duration) time.Duration {
	requireRandom(r, errNilRandom)
	if max <= 0 {
		return 0
	}

	bound := uint64(max) + 1
	return time.Duration(uint64(r.Int63()) % bound)
}

// randomOffsetInclusive returns a pseudo-random offset in [0, maxOffset].
//
// The helper is a semantic alias around randomDurationInclusive for callers that
// want the name to reflect offset arithmetic rather than absolute delay
// generation.
func randomOffsetInclusive(r RandomGenerator, maxOffset time.Duration) time.Duration {
	return randomDurationInclusive(r, maxOffset)
}
