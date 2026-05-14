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
	"time"

	"arcoris.dev/chrono/delay"
)

const (
	// errNilJitterSchedule is the stable diagnostic text used when a jitter
	// wrapper receives a nil child schedule.
	//
	// Jitter wrappers transform delay values produced by another schedule. They
	// cannot produce meaningful values without a child schedule to wrap. A nil
	// child indicates invalid package use rather than finite schedule exhaustion.
	errNilJitterSchedule = "jitter: nil schedule"

	// errNilJitterTransform is the stable diagnostic text used when the
	// package-local jitter wrapper is constructed with a nil transform function.
	//
	// Public callers do not provide jitter transforms directly. The diagnostic
	// exists for package-owned construction paths and tests so internal wrapper
	// wiring failures are reported at construction time.
	errNilJitterTransform = "jitter: nil transform"

	// errInvalidJitterFactor is the stable diagnostic text used when a jitter
	// algorithm receives an invalid factor.
	//
	// Factor-based jitter accepts finite non-negative values. Negative, NaN, and
	// infinite factors cannot describe a stable mechanical delay transform and are
	// rejected before a Schedule is created.
	errInvalidJitterFactor = "jitter: invalid factor"

	// errInvalidJitterRatio is the stable diagnostic text used when a jitter
	// algorithm receives an invalid ratio.
	//
	// Ratio-based jitter accepts finite values in [0, 1]. Values outside that
	// range can produce negative lower bounds or unbounded expansion, so they are
	// rejected at the constructor boundary.
	errInvalidJitterRatio = "jitter: invalid ratio"

	// errJitterScheduleReturnedNilSequence is the stable diagnostic text used
	// when the child schedule wrapped by a jitter wrapper returns a nil Sequence.
	//
	// Schedule.NewSequence must return a usable Sequence. Returning nil violates
	// the Schedule contract and would otherwise move a construction-time
	// programming error into a later retry, polling, reconnect, or cooldown path.
	errJitterScheduleReturnedNilSequence = "jitter: schedule returned nil Sequence"

	// errJitterScheduleReturnedNegativeDelay is the stable diagnostic text used
	// when a child sequence wrapped by a jitter wrapper returns a negative delay
	// with ok=true.
	//
	// Sequence implementations must return non-negative delays when ok=true.
	// Jitter wrappers do not repair negative child values because doing so would
	// hide a child Sequence contract violation.
	errJitterScheduleReturnedNegativeDelay = "jitter: schedule returned negative delay"

	// errJitterTransformReturnedNegativeDelay is the stable diagnostic text used
	// when a package-local jitter transform returns a negative delay.
	//
	// Concrete jitter algorithms must return non-negative delays. A negative
	// transformed delay indicates an internal package bug or an invalid custom
	// transform used by a package-local test.
	errJitterTransformReturnedNegativeDelay = "jitter: transform returned negative delay"
)

// jitterTransform applies one concrete jitter algorithm to a base delay.
//
// The base delay is always non-negative because jitterSequence validates child
// output before calling the transform. The random generator belongs to exactly
// one jitter sequence and is used only for inclusive duration offset draws.
//
// A transform must return a non-negative duration. It must not sleep, create
// timers, observe contexts, classify errors, execute operations, log, trace,
// export metrics, rate limit callers, schedule queue items, or make domain
// decisions.
type jitterTransform func(base time.Duration, r RandomGenerator) time.Duration

// newJitterSchedule returns a wrapper schedule that transforms child delays.
//
// The helper owns common option handling for concrete jitter constructors. It
// applies RandomOption values, obtains the RandomSource stored by the reusable
// Schedule, and delegates the final validated construction to
// newJitterScheduleWithSource.
func newJitterSchedule(sched delay.Schedule, xf jitterTransform, opts ...RandomOption) delay.Schedule {
	config := randomOptionsOf(opts...)
	return newJitterScheduleWithSource(sched, xf, config.source)
}

// newJitterScheduleWithSource wires shared jitter mechanics to source.
//
// The helper is package-local so tests can verify nil-source behavior directly
// without relying on RandomOption plumbing. The returned Schedule stores source;
// each Sequence created from it receives a random generator from source.NewRandom.
func newJitterScheduleWithSource(
	sched delay.Schedule,
	xf jitterTransform,
	src RandomSource,
) delay.Schedule {
	requireSchedule(sched, errNilJitterSchedule)
	if xf == nil {
		panic(errNilJitterTransform)
	}
	requireRandomSource(src, errNilRandomSource)

	return jitterSchedule{
		schedule:  sched,
		transform: xf,
		source:    src,
	}
}

// jitterSchedule is the reusable recipe behind concrete jitter wrappers.
//
// The type stores a child schedule, an algorithm transform, and a RandomSource.
// It does not store per-owner cursor state, previous delay state, timer state,
// retry state, context state, or mutable random generator state. Each
// NewSequence call creates a fresh child sequence and a fresh random generator
// from the source.
type jitterSchedule struct {
	// schedule is the child schedule whose available delays are transformed.
	//
	// The value is non-nil. newJitterScheduleWithSource validates it before a
	// jitterSchedule can be constructed.
	schedule delay.Schedule

	// transform is the concrete jitter algorithm applied to each child delay.
	//
	// The value is non-nil. It receives only non-negative child delays and must
	// return non-negative transformed delays.
	transform jitterTransform

	// source creates the per-sequence random generator used by transform.
	//
	// The value is non-nil. It is stored on the reusable schedule, while the
	// mutable generator returned by NewRandom belongs to one jitterSequence.
	source RandomSource
}

// NewSequence returns an independent jittered delay sequence.
//
// The method creates a fresh child sequence and a fresh random generator. It
// validates both boundary values immediately so nil child sequences and nil
// random generators are reported before the owner starts using the returned
// sequence.
func (s jitterSchedule) NewSequence() delay.Sequence {
	child := s.schedule.NewSequence()
	requireSequence(child, errJitterScheduleReturnedNilSequence)

	random := s.source.NewRandom()
	requireRandom(random, errNilRandom)

	return &jitterSequence{
		child:     child,
		transform: s.transform,
		random:    random,
	}
}

// jitterSequence is the per-owner delay stream produced by a jitter wrapper.
//
// The sequence delegates availability to child and applies transform only to
// available non-negative child delays. Exhaustion passes through unchanged. The
// random generator is sequence-owned mutable state and should not be shared with
// unrelated runtime owners.
type jitterSequence struct {
	// child is the wrapped per-owner delay sequence.
	//
	// The value is non-nil. jitterSchedule.NewSequence validates this before
	// constructing jitterSequence.
	child delay.Sequence

	// transform is the concrete jitter algorithm applied to available child
	// delays.
	transform jitterTransform

	// random is the sequence-owned pseudo-random generator used by transform.
	random RandomGenerator
}

// Next returns the next child delay after applying the configured transform.
//
// If the child sequence is exhausted, Next returns delay=0, ok=false without
// calling the transform. If the child returns a negative delay with ok=true, Next
// panics because the child has violated the Sequence contract. If the transform
// returns a negative delay, Next panics because jitter algorithms must preserve
// the non-negative delay contract.
func (s *jitterSequence) Next() (time.Duration, bool) {
	base, ok := s.child.Next()
	if !ok {
		return 0, false
	}
	requireNonNegativeSequenceDelay(base, ok, errJitterScheduleReturnedNegativeDelay)

	d := s.transform(base, s.random)
	if d < 0 {
		panic(errJitterTransformReturnedNegativeDelay)
	}

	return d, true
}
