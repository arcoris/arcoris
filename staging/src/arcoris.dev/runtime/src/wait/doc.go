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

// Package wait provides low-level context-aware runtime wait mechanics.
//
// The package owns mechanical waiting only: cancellable delays, owner-controlled
// timers, fixed-interval condition loops, deterministic condition composition,
// and wait-owned interruption/timeout errors. It does not own retry policy,
// backoff, randomized jitter, reusable delay schedules, scheduling policy,
// queues, health, logging, metrics, tracing, or lifecycle state.
//
// Public wait APIs panic on nil contexts, nil conditions, nil predicates, nil
// timers, and invalid fixed intervals. Callers that want an unconstrained
// context must pass context.Background explicitly.
//
// Until evaluates its condition immediately before sleeping and then repeats on
// the fixed interval. Condition-owned errors are returned unchanged. Context
// stops owned by wait primitives are classified as ErrInterrupted or ErrTimeout
// while preserving cancellation causes where the standard context package
// exposes them.
//
// Timer values must be created with NewTimer. The zero value is invalid and
// Timer follows a single-owner coordination model; callers sharing Stop,
// StopAndDrain, Reset, Wait, or C across goroutines must synchronize that
// ownership themselves.
//
// Production code in this package depends only on the Go standard library.
package wait
