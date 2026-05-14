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

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestSatisfiedReturnsDone(t *testing.T) {
	t.Parallel()

	done, err := Satisfied(context.Background())

	if !done {
		t.Fatal("Satisfied(...) done = false, want true")
	}
	if err != nil {
		t.Fatalf("Satisfied(...) err = %v, want nil", err)
	}
}

func TestSatisfiedIgnoresStoppedContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done, err := Satisfied(ctx)

	if !done {
		t.Fatal("Satisfied(cancelled context) done = false, want true")
	}
	if err != nil {
		t.Fatalf("Satisfied(cancelled context) err = %v, want nil", err)
	}
}

func TestUnsatisfiedReturnsContinue(t *testing.T) {
	t.Parallel()

	done, err := Unsatisfied(context.Background())

	if done {
		t.Fatal("Unsatisfied(...) done = true, want false")
	}
	if err != nil {
		t.Fatalf("Unsatisfied(...) err = %v, want nil", err)
	}
}

func TestUnsatisfiedIgnoresStoppedContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done, err := Unsatisfied(ctx)

	if done {
		t.Fatal("Unsatisfied(cancelled context) done = true, want false")
	}
	if err != nil {
		t.Fatalf("Unsatisfied(cancelled context) err = %v, want nil", err)
	}
}

func TestPredicatePanicsOnNilPredicate(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilPredicate, func() {
		_ = Predicate(nil)
	})
}

func TestPredicatePassesExactContext(t *testing.T) {
	t.Parallel()

	type key struct{}

	ctx := context.WithValue(context.Background(), key{}, "value")
	var got context.Context
	condition := Predicate(func(ctx context.Context) bool {
		got = ctx
		return true
	})

	done, err := condition(ctx)

	if err != nil {
		t.Fatalf("Predicate(...) err = %v, want nil", err)
	}
	if !done {
		t.Fatal("Predicate(...) done = false, want true")
	}
	if got != ctx {
		t.Fatal("predicate received different context, want exact context")
	}
}

func TestPredicateReturnsConditionState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want bool
	}{
		{
			name: "true",
			want: true,
		},
		{
			name: "false",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			condition := Predicate(func(context.Context) bool {
				return tc.want
			})

			done, err := condition(context.Background())

			if done != tc.want {
				t.Fatalf("Predicate(...) done = %v, want %v", done, tc.want)
			}
			if err != nil {
				t.Fatalf("Predicate(...) err = %v, want nil", err)
			}
		})
	}
}

func TestPredicateDoesNotRecoverPanic(t *testing.T) {
	t.Parallel()

	panicValue := "predicate panic"
	condition := Predicate(func(context.Context) bool {
		panic(panicValue)
	})

	mustPanicWith(t, panicValue, func() {
		_, _ = condition(context.Background())
	})
}

func TestNotPanicsOnNilCondition(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilCondition, func() {
		_ = Not(nil)
	})
}

func TestNotInvertsSatisfiedState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   bool
		want bool
	}{
		{
			name: "true becomes false",
			in:   true,
			want: false,
		},
		{
			name: "false becomes true",
			in:   false,
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			condition := Not(func(context.Context) (bool, error) {
				return tc.in, nil
			})

			done, err := condition(context.Background())

			if done != tc.want {
				t.Fatalf("Not(...) done = %v, want %v", done, tc.want)
			}
			if err != nil {
				t.Fatalf("Not(...) err = %v, want nil", err)
			}
		})
	}
}

func TestNotReturnsErrorUnchanged(t *testing.T) {
	t.Parallel()

	conditionErr := errors.New("condition failed")
	condition := Not(func(context.Context) (bool, error) {
		return true, conditionErr
	})

	done, err := condition(context.Background())

	if done {
		t.Fatal("Not(...) done = true after error, want false")
	}
	if err != conditionErr {
		t.Fatalf("Not(...) err = %v, want exact error %v", err, conditionErr)
	}
}

func TestNotEvaluatesConditionOnce(t *testing.T) {
	t.Parallel()

	calls := 0
	condition := Not(func(context.Context) (bool, error) {
		calls++
		return false, nil
	})

	_, err := condition(context.Background())

	if err != nil {
		t.Fatalf("Not(...) err = %v, want nil", err)
	}
	if calls != 1 {
		t.Fatalf("condition calls = %d, want 1", calls)
	}
}

func TestAllPanicsOnNilCondition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func()
	}{
		{
			name: "first",
			call: func() { _ = All(nil) },
		},
		{
			name: "rest",
			call: func() { _ = All(Satisfied, nil) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilCondition, tc.call)
		})
	}
}

func TestAllEvaluatesConditionsSequentially(t *testing.T) {
	t.Parallel()

	var order []string
	condition := All(
		recordCondition(&order, "first", true, nil),
		recordCondition(&order, "second", true, nil),
		recordCondition(&order, "third", true, nil),
	)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("All(...) err = %v, want nil", err)
	}
	if !done {
		t.Fatal("All(...) done = false, want true")
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second", "third"})
}

func TestAllShortCircuitsOnFirstUnsatisfiedCondition(t *testing.T) {
	t.Parallel()

	var order []string
	condition := All(
		recordCondition(&order, "first", true, nil),
		recordCondition(&order, "second", false, nil),
		recordCondition(&order, "third", true, nil),
	)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("All(...) err = %v, want nil", err)
	}
	if done {
		t.Fatal("All(...) done = true, want false")
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second"})
}

func TestAllShortCircuitsOnFirstError(t *testing.T) {
	t.Parallel()

	conditionErr := errors.New("condition failed")
	var order []string
	condition := All(
		recordCondition(&order, "first", true, nil),
		recordCondition(&order, "second", true, conditionErr),
		recordCondition(&order, "third", true, nil),
	)

	done, err := condition(context.Background())

	if done {
		t.Fatal("All(...) done = true after error, want false")
	}
	if err != conditionErr {
		t.Fatalf("All(...) err = %v, want exact error %v", err, conditionErr)
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second"})
}

func TestAllReturnsTrueOnlyWhenAllSatisfied(t *testing.T) {
	t.Parallel()

	condition := All(Satisfied, Satisfied)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("All(...) err = %v, want nil", err)
	}
	if !done {
		t.Fatal("All(...) done = false, want true")
	}
}

func TestAllSnapshotsConditionList(t *testing.T) {
	t.Parallel()

	var order []string
	rest := []ConditionFunc{recordCondition(&order, "second", true, nil)}
	condition := All(recordCondition(&order, "first", true, nil), rest...)
	rest[0] = recordCondition(&order, "mutated", false, nil)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("All(...) err = %v, want nil", err)
	}
	if !done {
		t.Fatal("All(...) done = false, want true")
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second"})
}

func TestAnyPanicsOnNilCondition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func()
	}{
		{
			name: "first",
			call: func() { _ = Any(nil) },
		},
		{
			name: "rest",
			call: func() { _ = Any(Unsatisfied, nil) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilCondition, tc.call)
		})
	}
}

func TestAnyEvaluatesConditionsSequentially(t *testing.T) {
	t.Parallel()

	var order []string
	condition := Any(
		recordCondition(&order, "first", false, nil),
		recordCondition(&order, "second", false, nil),
		recordCondition(&order, "third", false, nil),
	)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("Any(...) err = %v, want nil", err)
	}
	if done {
		t.Fatal("Any(...) done = true, want false")
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second", "third"})
}

func TestAnyShortCircuitsOnFirstSatisfiedCondition(t *testing.T) {
	t.Parallel()

	var order []string
	condition := Any(
		recordCondition(&order, "first", false, nil),
		recordCondition(&order, "second", true, nil),
		recordCondition(&order, "third", true, nil),
	)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("Any(...) err = %v, want nil", err)
	}
	if !done {
		t.Fatal("Any(...) done = false, want true")
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second"})
}

func TestAnyReturnsErrorBeforeSuccess(t *testing.T) {
	t.Parallel()

	conditionErr := errors.New("condition failed")
	var order []string
	condition := Any(
		recordCondition(&order, "first", false, nil),
		recordCondition(&order, "second", false, conditionErr),
		recordCondition(&order, "third", true, nil),
	)

	done, err := condition(context.Background())

	if done {
		t.Fatal("Any(...) done = true after error, want false")
	}
	if err != conditionErr {
		t.Fatalf("Any(...) err = %v, want exact error %v", err, conditionErr)
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second"})
}

func TestAnyReturnsFalseIfNoConditionSatisfied(t *testing.T) {
	t.Parallel()

	condition := Any(Unsatisfied, Unsatisfied)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("Any(...) err = %v, want nil", err)
	}
	if done {
		t.Fatal("Any(...) done = true, want false")
	}
}

func TestAnySnapshotsConditionList(t *testing.T) {
	t.Parallel()

	var order []string
	rest := []ConditionFunc{recordCondition(&order, "second", false, nil)}
	condition := Any(recordCondition(&order, "first", false, nil), rest...)
	rest[0] = recordCondition(&order, "mutated", true, nil)

	done, err := condition(context.Background())

	if err != nil {
		t.Fatalf("Any(...) err = %v, want nil", err)
	}
	if done {
		t.Fatal("Any(...) done = true, want false")
	}
	mustEqualStringSlice(t, "condition order", order, []string{"first", "second"})
}

func recordCondition(order *[]string, name string, done bool, err error) ConditionFunc {
	return func(context.Context) (bool, error) {
		*order = append(*order, name)
		return done, err
	}
}

func mustEqualStringSlice(t *testing.T, name string, got []string, want []string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}
