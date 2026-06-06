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

// entry builds a valid test entry.
func entry(resource string, amount uint64) capacity.Entry {
	return capacity.Entry{
		Resource: capacity.MustResource(resource),
		Amount:   capacity.Amount(amount),
	}
}

// vector builds a valid test vector.
func vector(t testing.TB, entries ...capacity.Entry) capacity.Vector {
	t.Helper()

	vector, err := capacity.NewVector(entries...)
	if err != nil {
		t.Fatalf("NewVector() error = %v", err)
	}

	return vector
}

// demand builds a valid test demand.
func demand(t testing.TB, entries ...capacity.Entry) capacity.Demand {
	t.Helper()

	demand, err := capacity.NewDemand(entries...)
	if err != nil {
		t.Fatalf("NewDemand() error = %v", err)
	}

	return demand
}

// requireVector checks a vector's canonical entries.
func requireVector(t testing.TB, got capacity.Vector, want ...capacity.Entry) {
	t.Helper()

	if !reflect.DeepEqual(got.Entries(), want) {
		t.Fatalf("vector entries = %#v, want %#v", got.Entries(), want)
	}
}

// requirePanicIs checks that fn panics with a matching sentinel error.
func requirePanicIs(t *testing.T, sentinel error, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %v", sentinel)
		}

		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic = %T, want error", recovered)
		}
		if !errors.Is(err, sentinel) {
			t.Fatalf("panic error = %v, want errors.Is(..., %v)", err, sentinel)
		}
	}()

	fn()
}
