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

package delay

import "time"

const (
	// errNilChainSchedule is the stable diagnostic text used when Chain receives
	// a nil child schedule.
	//
	// Chain composes concrete child schedules in order. A nil child cannot create
	// a per-owner Sequence and is invalid package use rather than finite
	// exhaustion.
	errNilChainSchedule = "delay: nil chain schedule"

	// errChainScheduleReturnedNilSequence is the stable diagnostic text used when
	// a child schedule wrapped by Chain returns a nil Sequence.
	//
	// Chain creates child sequences at NewSequence time so each owner receives an
	// independent stream. A nil child sequence violates the Schedule contract at
	// that boundary.
	errChainScheduleReturnedNilSequence = "delay: chain schedule returned nil Sequence"

	// errChainScheduleReturnedNegativeDelay is the stable diagnostic text used
	// when a child sequence wrapped by Chain returns a negative delay with
	// ok=true.
	//
	// Chain forwards available child delays unchanged, so it must reject negative
	// available values rather than hiding a child Sequence contract violation.
	errChainScheduleReturnedNegativeDelay = "delay: chain schedule returned negative delay"
)

// Chain returns a schedule that consumes child schedules sequentially.
//
// Each sequence created by the returned Schedule starts with the first child
// schedule. When the current child sequence reports exhaustion with ok=false,
// the chain advances to the next child sequence. If a child sequence is
// infinite, later children are never observed. An empty chain is valid and
// produces an immediately exhausted sequence.
//
// Chain is useful for deterministic profiles that combine a finite prefix with
// a growing or infinite tail, for example:
//
//	delay.Chain(
//	    delay.Delays(0),
//	    delay.Exponential(100*time.Millisecond, 2),
//	)
//
// That schedule produces an immediate first value, then continues with the
// exponential sequence.
//
// Chain preserves zero delays as available values. A child result of delay=0,
// ok=true is forwarded as immediate continuation and is not confused with
// exhaustion. A child result with ok=false advances the chain and the returned
// delay is ignored, even if that ignored value is negative.
//
// The returned Schedule copies caller-owned slice state. Mutating the caller's
// slice after construction does not change the chain. Each call to NewSequence
// creates independent child sequences.
//
// Chain does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, randomize delays, log, trace,
// export metrics, rate limit callers, schedule queue items, or make domain
// decisions.
//
// Chain panics when any child schedule is nil. NewSequence panics if any child
// schedule returns a nil Sequence. Next panics if a child sequence returns a
// negative delay with ok=true.
func Chain(schedules ...Schedule) Schedule {
	switch len(schedules) {
	case 0:
		return emptyChainSchedule{}
	case 1:
		child := schedules[0]
		requireSchedule(child, errNilChainSchedule)

		return singleChainSchedule{child: child}
	case 2:
		first := schedules[0]
		second := schedules[1]
		requireSchedule(first, errNilChainSchedule)
		requireSchedule(second, errNilChainSchedule)

		return pairChainSchedule{
			first:  first,
			second: second,
		}
	default:
		copied := make([]Schedule, len(schedules))
		copy(copied, schedules)

		for _, schedule := range copied {
			requireSchedule(schedule, errNilChainSchedule)
		}

		return chainSchedule{schedules: copied}
	}
}

// Implementation note: Chain uses separate shapes for small chains.
//
// The common composition is a finite prefix followed by an infinite tail:
// Chain(Delays(...), Exponential(...)) or Chain(Delays(...), Fixed(...)).
// Keeping the zero-, one-, and two-child cases inline avoids heap slices on that
// common path. Larger chains still copy caller-owned slice state to preserve the
// immutable Schedule contract.

// emptyChainSchedule is the reusable recipe behind Chain with no children.
type emptyChainSchedule struct{}

// NewSequence returns an already exhausted sequence.
func (emptyChainSchedule) NewSequence() Sequence {
	return exhaustedSequence{}
}

// singleChainSchedule is the reusable recipe behind a one-child Chain.
//
// It still wraps the child sequence so Chain can enforce its boundary checks and
// stop calling the child after exhaustion.
type singleChainSchedule struct {
	// child is the only child schedule.
	child Schedule
}

// NewSequence creates an independent single-child chain sequence.
func (s singleChainSchedule) NewSequence() Sequence {
	child := s.child.NewSequence()
	requireSequence(child, errChainScheduleReturnedNilSequence)

	return &singleChainSequence{child: child}
}

// pairChainSchedule is the reusable recipe behind the common two-child Chain.
type pairChainSchedule struct {
	// first is the first child schedule.
	first Schedule

	// second is the second child schedule.
	second Schedule
}

// NewSequence creates an independent two-child chain sequence.
func (s pairChainSchedule) NewSequence() Sequence {
	first := s.first.NewSequence()
	requireSequence(first, errChainScheduleReturnedNilSequence)

	second := s.second.NewSequence()
	requireSequence(second, errChainScheduleReturnedNilSequence)

	return &pairChainSequence{
		first:  first,
		second: second,
	}
}

// chainSchedule is the reusable recipe behind a Chain with three or more
// children.
//
// The type stores a copied child schedule slice and no iteration state.
// Per-owner child sequence state is held by chainSequence values created by
// NewSequence.
type chainSchedule struct {
	// schedules is the copied ordered list of child schedules.
	//
	// The slice MUST NOT be mutated after construction. The public constructor
	// copies caller input before storing it.
	schedules []Schedule
}

// NewSequence creates an independent chained sequence.
//
// All child sequences are created up front so child Schedule contract violations
// are reported at sequence creation time, before a runtime owner starts
// consuming delay values.
func (s chainSchedule) NewSequence() Sequence {
	sequences := make([]Sequence, len(s.schedules))
	for i, schedule := range s.schedules {
		sequence := schedule.NewSequence()
		requireSequence(sequence, errChainScheduleReturnedNilSequence)

		sequences[i] = sequence
	}

	return &chainSequence{sequences: sequences}
}

// singleChainSequence is the per-owner delay stream for a one-child Chain.
//
// The sequence stops calling its child after the child has reported exhaustion.
type singleChainSequence struct {
	// child is the wrapped per-owner sequence.
	child Sequence

	// exhausted records that child has already returned ok=false.
	exhausted bool
}

// Next returns the next child delay or reports exhaustion.
func (s *singleChainSequence) Next() (time.Duration, bool) {
	if s.exhausted {
		return 0, false
	}

	d, ok := s.child.Next()
	if !ok {
		s.exhausted = true
		return 0, false
	}
	requireNonNegativeSequenceDelay(d, ok, errChainScheduleReturnedNegativeDelay)

	return d, true
}

// pairChainSequence is the per-owner delay stream for a two-child Chain.
type pairChainSequence struct {
	// first is the first child sequence.
	first Sequence

	// second is the second child sequence.
	second Sequence

	// next is the current child index.
	next int
}

// Next returns the next available child delay or reports exhaustion.
func (s *pairChainSequence) Next() (time.Duration, bool) {
	if s.next == 0 {
		d, ok := s.first.Next()
		requireNonNegativeSequenceDelay(d, ok, errChainScheduleReturnedNegativeDelay)
		if ok {
			return d, true
		}
		s.next = 1
	}

	if s.next == 1 {
		d, ok := s.second.Next()
		requireNonNegativeSequenceDelay(d, ok, errChainScheduleReturnedNegativeDelay)
		if ok {
			return d, true
		}
		s.next = 2
	}

	return 0, false
}

// chainSequence is the per-owner delay stream produced by large Chain values.
//
// The sequence owns a cursor into its child sequence list. Once a child sequence
// reports exhaustion, the cursor advances permanently and that child is not
// called again.
type chainSequence struct {
	// sequences is the ordered list of child sequences created for this owner.
	sequences []Sequence

	// next is the index of the current child sequence.
	next int
}

// Next returns the next available child delay or reports exhaustion.
//
// Exhausted child sequences are skipped. Available child delays are forwarded
// unchanged after the non-negative Sequence contract check.
func (s *chainSequence) Next() (time.Duration, bool) {
	for s.next < len(s.sequences) {
		d, ok := s.sequences[s.next].Next()
		if !ok {
			s.next++
			continue
		}
		requireNonNegativeSequenceDelay(d, ok, errChainScheduleReturnedNegativeDelay)

		return d, true
	}

	return 0, false
}

// exhaustedSequence is an always-exhausted sequence used by empty Chain.
type exhaustedSequence struct{}

// Next reports exhaustion.
func (exhaustedSequence) Next() (time.Duration, bool) {
	return 0, false
}
