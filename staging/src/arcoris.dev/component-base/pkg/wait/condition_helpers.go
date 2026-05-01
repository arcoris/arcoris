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

package wait

import "context"

const (
	// errNilPredicate is the panic value used when a helper receives a nil
	// predicate function.
	//
	// A nil predicate is a construction-time programming error. Delaying the
	// failure until the condition is evaluated would make the wait loop fail far
	// away from the invalid helper call.
	errNilPredicate = "wait: nil predicate"

	// errNilCondition is the panic value used when a helper receives a nil
	// condition function.
	//
	// A nil condition is a construction-time programming error. Wait helpers fail
	// fast so invalid condition graphs are detected before a runtime wait loop
	// starts.
	errNilCondition = "wait: nil condition"
)

// Satisfied is a condition that is already complete.
//
// Satisfied is useful when a caller needs an explicit ConditionFunc value for an
// already-ready state, a test fixture, a disabled wait branch, or a composition
// path that should not delay the owning wait operation.
//
// The supplied context is intentionally ignored. Cancellation, timeout, and
// shutdown policy belong to the wait operation that evaluates the condition, not
// to this constant helper condition.
func Satisfied(context.Context) (done bool, err error) {
	return true, nil
}

// Unsatisfied is a condition that is not yet complete.
//
// Unsatisfied is useful when a caller needs an explicit ConditionFunc value for
// a pending state, a test fixture, a permanently-blocked branch, or a
// composition path that should keep the owning wait operation running until some
// external policy stops it.
//
// The supplied context is intentionally ignored. Cancellation, timeout, and
// shutdown policy belong to the wait operation that evaluates the condition, not
// to this constant helper condition.
func Unsatisfied(context.Context) (done bool, err error) {
	return false, nil
}

// Predicate adapts a context-aware boolean predicate into a ConditionFunc.
//
// Predicate is for conditions whose evaluation cannot fail and only needs to
// report whether the wait should stop. The predicate receives the same context
// that the wait operation passes to the resulting condition.
//
// Predicate deliberately accepts func(context.Context) bool, not func() bool.
// Context-free predicate adapters are avoided so new code does not grow a
// parallel context-free condition model. A predicate that does not need the
// context can ignore it explicitly at the call site.
//
// Predicate does not recover panics raised by predicate. Panic recovery policy,
// if required, belongs to the wait loop owner or to an explicit higher-level
// wrapper.
//
// Predicate panics when predicate is nil.
func Predicate(predicate func(context.Context) bool) ConditionFunc {
	requirePredicate(predicate)

	return func(ctx context.Context) (done bool, err error) {
		return predicate(ctx), nil
	}
}

// Not returns a condition that inverts the completion result of condition.
//
// If condition returns done=true with a nil error, the returned condition reports
// done=false. If condition returns done=false with a nil error, the returned
// condition reports done=true.
//
// If condition returns an error, Not returns that error unchanged and reports the
// condition as not completed. Errors are terminal wait results and are not
// inverted.
//
// Not evaluates condition at most once per invocation and passes the received
// context through unchanged.
//
// Not panics when condition is nil.
func Not(condition ConditionFunc) ConditionFunc {
	requireCondition(condition)

	return func(ctx context.Context) (done bool, err error) {
		done, err = condition(ctx)
		if err != nil {
			return false, err
		}

		return !done, nil
	}
}

// All returns a condition that is satisfied only when every supplied condition is
// satisfied.
//
// Conditions are evaluated sequentially in the order they are provided. All
// short-circuits on the first unsatisfied condition and returns done=false with a
// nil error. If a condition returns an error, All stops immediately and returns
// that error unchanged.
//
// All is a deterministic sequential conjunction. It does not evaluate
// conditions concurrently, does not aggregate errors, does not retry failed
// conditions, and does not recover panics raised by condition implementations.
//
// All snapshots the supplied condition list at construction time. Mutating the
// caller-owned variadic slice after All returns cannot change the composed
// condition graph.
//
// All requires at least one condition by construction. It panics when any
// supplied condition is nil.
func All(first ConditionFunc, rest ...ConditionFunc) ConditionFunc {
	conditions := conditionsOf(first, rest...)

	return func(ctx context.Context) (done bool, err error) {
		for _, condition := range conditions {
			done, err = condition(ctx)
			if err != nil {
				return false, err
			}
			if !done {
				return false, nil
			}
		}

		return true, nil
	}
}

// Any returns a condition that is satisfied when at least one supplied condition
// is satisfied.
//
// Conditions are evaluated sequentially in the order they are provided. Any
// short-circuits on the first satisfied condition and returns done=true with a
// nil error. If a condition returns an error before any condition is satisfied,
// Any stops immediately and returns that error unchanged.
//
// Any is a deterministic sequential disjunction. It does not evaluate conditions
// concurrently, does not suppress errors that occur before success, does not
// retry failed conditions, and does not recover panics raised by condition
// implementations.
//
// Any snapshots the supplied condition list at construction time. Mutating the
// caller-owned variadic slice after Any returns cannot change the composed
// condition graph.
//
// Any requires at least one condition by construction. It panics when any
// supplied condition is nil.
func Any(first ConditionFunc, rest ...ConditionFunc) ConditionFunc {
	conditions := conditionsOf(first, rest...)

	return func(ctx context.Context) (done bool, err error) {
		for _, condition := range conditions {
			done, err = condition(ctx)
			if err != nil {
				return false, err
			}
			if done {
				return true, nil
			}
		}

		return false, nil
	}
}

// conditionsOf validates and snapshots a condition list.
//
// The returned slice is independent from the caller-owned variadic slice so a
// caller cannot mutate the composition after the helper returns.
//
// The first argument is separate from rest so public composition helpers require
// at least one condition by construction and do not need to define surprising
// empty-list semantics.
func conditionsOf(first ConditionFunc, rest ...ConditionFunc) []ConditionFunc {
	requireCondition(first)

	conditions := make([]ConditionFunc, 0, 1+len(rest))
	conditions = append(conditions, first)

	for _, condition := range rest {
		requireCondition(condition)
		conditions = append(conditions, condition)
	}

	return conditions
}

// requirePredicate panics when predicate is nil.
//
// Nil predicates are rejected at helper construction time instead of being
// allowed to fail later inside a wait loop.
func requirePredicate(predicate func(context.Context) bool) {
	if predicate == nil {
		panic(errNilPredicate)
	}
}

// requireCondition panics when condition is nil.
//
// Nil conditions are rejected at helper construction time instead of being
// allowed to fail later inside a wait loop.
func requireCondition(condition ConditionFunc) {
	if condition == nil {
		panic(errNilCondition)
	}
}
