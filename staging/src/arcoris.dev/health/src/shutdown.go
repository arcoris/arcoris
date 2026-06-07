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

package health

import (
	"context"
)

// NewShutdownCheck returns a Checker that reports unhealthy after done is closed.
//
// The returned checker is intended for readiness-style registration. A graceful
// shutdown should normally stop new work before it affects liveness. Registering
// this checker under TargetLive is possible, but it changes shutdown into a
// liveness failure and should be an explicit owner decision.
//
// While done is open, the checker reports StatusHealthy. After done is closed,
// it reports StatusUnhealthy with ReasonShuttingDown.
//
// done MUST be non-nil.
func NewShutdownCheck(name string, done <-chan struct{}) (Checker, error) {
	return newSourceChannelCheck(
		name,
		done,
		ReasonShuttingDown,
		"component is shutting down",
	)
}

// MustShutdownCheck returns a shutdown checker and panics if construction fails.
//
// MustShutdownCheck is intended for package-level declarations and tests where
// invalid wiring is a programmer error. Runtime configuration paths SHOULD use
// NewShutdownCheck and return the error to the owner.
func MustShutdownCheck(name string, done <-chan struct{}) Checker {
	checker, err := NewShutdownCheck(name, done)
	if err != nil {
		panic(err)
	}

	return checker
}

// NewDrainCheck returns a Checker that reports unhealthy after draining is
// closed.
//
// The returned checker is intended for readiness or admission registration.
// Draining means the component is intentionally refusing new work while it
// finishes in-flight work. It should normally affect readiness, not liveness.
//
// While draining is open, the checker reports StatusHealthy. After draining is
// closed, it reports StatusUnhealthy with ReasonDraining.
//
// draining MUST be non-nil.
func NewDrainCheck(name string, draining <-chan struct{}) (Checker, error) {
	return newSourceChannelCheck(
		name,
		draining,
		ReasonDraining,
		"component is draining",
	)
}

// MustDrainCheck returns a drain checker and panics if construction fails.
//
// MustDrainCheck is intended for package-level declarations and tests where
// invalid wiring is a programmer error. Runtime configuration paths SHOULD use
// NewDrainCheck and return the error to the owner.
func MustDrainCheck(name string, draining <-chan struct{}) Checker {
	checker, err := NewDrainCheck(name, draining)
	if err != nil {
		panic(err)
	}

	return checker
}

// NewContextShutdownCheck returns a Checker that reports unhealthy after source
// is cancelled.
//
// The returned checker is useful when shutdown state is already represented by a
// context, for example by a run owner, lifecycle owner, signal controller, or
// parent component. The checker observes source, not the evaluation context
// passed to Check.
//
// While source is not cancelled, the checker reports StatusHealthy. After source
// is cancelled, it reports StatusUnhealthy with ReasonShuttingDown and preserves
// context.Cause(source) as Result.Cause when a cause is available.
//
// source MUST be non-nil.
func NewContextShutdownCheck(name string, source context.Context) (Checker, error) {
	return newSourceContextCheck(
		name,
		source,
		ReasonShuttingDown,
		"component is shutting down",
	)
}

// MustContextShutdownCheck returns a context-backed shutdown checker and panics
// if construction fails.
//
// MustContextShutdownCheck is intended for package-level declarations and tests
// where invalid wiring is a programmer error. Runtime configuration paths SHOULD
// use NewContextShutdownCheck and return the error to the owner.
func MustContextShutdownCheck(name string, source context.Context) Checker {
	checker, err := NewContextShutdownCheck(name, source)
	if err != nil {
		panic(err)
	}

	return checker
}

// NewContextDrainCheck returns a Checker that reports unhealthy after source is
// cancelled.
//
// The returned checker is useful when drain state is represented by a context
// rather than a bare channel. It observes source, not the evaluation context
// passed to Check.
//
// While source is not cancelled, the checker reports StatusHealthy. After source
// is cancelled, it reports StatusUnhealthy with ReasonDraining and preserves
// context.Cause(source) as Result.Cause when a cause is available.
//
// source MUST be non-nil.
func NewContextDrainCheck(name string, source context.Context) (Checker, error) {
	return newSourceContextCheck(
		name,
		source,
		ReasonDraining,
		"component is draining",
	)
}

// MustContextDrainCheck returns a context-backed drain checker and panics if
// construction fails.
//
// MustContextDrainCheck is intended for package-level declarations and tests
// where invalid wiring is a programmer error. Runtime configuration paths SHOULD
// use NewContextDrainCheck and return the error to the owner.
func MustContextDrainCheck(name string, source context.Context) Checker {
	checker, err := NewContextDrainCheck(name, source)
	if err != nil {
		panic(err)
	}

	return checker
}
