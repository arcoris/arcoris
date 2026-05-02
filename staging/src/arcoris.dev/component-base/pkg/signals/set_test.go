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

package signals

import (
	"os"
	"strings"
	"testing"
)

type nonComparableSignal struct {
	name string
	data []int
}

func (s nonComparableSignal) Signal() {}

func (s nonComparableSignal) String() string { return s.name }

func TestClonePreservesNilAndCopiesNonNil(t *testing.T) {
	t.Parallel()

	if Clone(nil) != nil {
		t.Fatal("Clone(nil) returned non-nil")
	}

	original := []os.Signal{testSIGINT}
	clone := Clone(original)
	clone[0] = testSIGTERM

	if !sameSignal(original[0], testSIGINT) {
		t.Fatal("Clone did not isolate the returned slice")
	}
}

func TestUniquePreservesOrderAndRemovesDuplicates(t *testing.T) {
	t.Parallel()

	got := Unique([]os.Signal{testSIGINT, testSIGTERM, testSIGINT, testSIGHUP})
	want := []os.Signal{testSIGINT, testSIGTERM, testSIGHUP}

	assertSignalSlice(t, got, want)
}

func TestUniqueAndMergeReturnNilForEmptyInput(t *testing.T) {
	t.Parallel()

	if got := Unique(nil); got != nil {
		t.Fatalf("Unique(nil) = %v, want nil", got)
	}
	if got := Merge(); got != nil {
		t.Fatalf("Merge() = %v, want nil", got)
	}
}

func TestUniqueHandlesNonComparableSignals(t *testing.T) {
	t.Parallel()

	left := nonComparableSignal{name: "custom", data: []int{1}}
	right := nonComparableSignal{name: "custom", data: []int{2}}

	got := Unique([]os.Signal{left, right})

	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if !sameSignal(got[0], left) {
		t.Fatalf("signal = %v, want %v", got[0], left)
	}
}

func TestMergePreservesSetOrderAndUniqueness(t *testing.T) {
	t.Parallel()

	got := Merge(
		[]os.Signal{testSIGINT, testSIGTERM},
		[]os.Signal{testSIGTERM, testSIGHUP},
	)
	want := []os.Signal{testSIGINT, testSIGTERM, testSIGHUP}

	assertSignalSlice(t, got, want)
}

func TestContainsReportsSignalMembership(t *testing.T) {
	t.Parallel()

	set := []os.Signal{testSIGINT, testSIGTERM}
	if !Contains(set, testSIGTERM) {
		t.Fatal("Contains did not find present signal")
	}
	if Contains(set, testSIGHUP) {
		t.Fatal("Contains found missing signal")
	}
}

func TestSignalSetHelpersRejectNilSignals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func()
	}{
		{name: "unique nil member", fn: func() { Unique([]os.Signal{nil}) }},
		{name: "merge nil member", fn: func() { Merge([]os.Signal{testSIGINT}, []os.Signal{nil}) }},
		{name: "contains nil target", fn: func() { Contains([]os.Signal{testSIGINT}, nil) }},
		{name: "contains nil member", fn: func() { Contains([]os.Signal{nil}, testSIGINT) }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilSignalSetSignal, tt.fn)
		})
	}
}

func TestSignalSetOrShutdownSignalsDefaultsEmptyInput(t *testing.T) {
	t.Parallel()

	got := signalSetOrShutdownSignals(nil)

	if len(got) == 0 {
		t.Fatal("empty input did not default to shutdown signals")
	}
	got[0] = testSIGHUP
	if sameSignal(signalSetOrShutdownSignals(nil)[0], testSIGHUP) {
		t.Fatal("default shutdown set was mutated through returned slice")
	}
}

func TestSignalSetOrShutdownSignalsNormalizesExplicitInput(t *testing.T) {
	t.Parallel()

	got := signalSetOrShutdownSignals([]os.Signal{testSIGINT, testSIGINT, testSIGTERM})
	want := []os.Signal{testSIGINT, testSIGTERM}

	assertSignalSlice(t, got, want)
}

func TestSignalKeyUsesTypeAndString(t *testing.T) {
	t.Parallel()

	key := signalKey(testSIGINT)

	if !strings.Contains(key, "signals.testSignal") {
		t.Fatalf("signalKey(%v) = %q, want type name", testSIGINT, key)
	}
	if !strings.Contains(key, testSIGINT.String()) {
		t.Fatalf("signalKey(%v) = %q, want signal string", testSIGINT, key)
	}
}

func TestTypeNameHandlesNil(t *testing.T) {
	t.Parallel()

	if got := typeName(nil); got != "<nil>" {
		t.Fatalf("typeName(nil) = %q, want <nil>", got)
	}
}

func assertSignalSlice(t *testing.T, got, want []os.Signal) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if !sameSignal(got[i], want[i]) {
			t.Fatalf("signal[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
