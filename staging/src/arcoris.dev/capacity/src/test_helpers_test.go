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

package capacity_test

import (
	"errors"
	"reflect"
	"testing"

	"arcoris.dev/capacity"
)

// entry builds one valid test entry from compact literals.
func entry(resource string, amount uint64) capacity.Entry {
	return capacity.Entry{
		Resource: capacity.MustResource(resource),
		Amount:   capacity.Amount(amount),
	}
}

// vector builds a test vector and fails the current test on invalid fixtures.
func vector(t *testing.T, entries ...capacity.Entry) capacity.Vector {
	t.Helper()

	v, err := capacity.NewVector(entries...)
	if err != nil {
		t.Fatalf("NewVector() error = %v", err)
	}

	return v
}

// demand builds a test demand and fails the current test on invalid fixtures.
func demand(t *testing.T, entries ...capacity.Entry) capacity.Demand {
	t.Helper()

	d, err := capacity.NewDemand(entries...)
	if err != nil {
		t.Fatalf("NewDemand() error = %v", err)
	}

	return d
}

// requireEntries compares canonical entries without hiding ordering mistakes.
func requireEntries(t *testing.T, got []capacity.Entry, want ...capacity.Entry) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("entries = %#v, want %#v", got, want)
	}
}

// requireVector compares a vector against expected canonical entries.
func requireVector(t *testing.T, got capacity.Vector, want ...capacity.Entry) {
	t.Helper()

	requireEntries(t, got.Entries(), want...)
}

// requirePanicIs verifies that fn panics with an error wrapping sentinel.
func requirePanicIs(t *testing.T, sentinel error, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("function did not panic")
		}

		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic = %#v, want error", recovered)
		}

		if !errors.Is(err, sentinel) {
			t.Fatalf("panic error = %v, want errors.Is(..., %v)", err, sentinel)
		}
	}()

	fn()
}
