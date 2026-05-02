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
	"testing"
)

func TestSignalSetsReturnIndependentSlices(t *testing.T) {
	left := Shutdown()
	if len(left) == 0 {
		t.Fatal("Shutdown returned an empty set")
	}

	original := left[0]
	left[0] = testSIGHUP

	right := Shutdown()
	if !sameSignal(right[0], original) {
		t.Fatalf("Shutdown was mutated through returned slice: got %v, want %v", right[0], original)
	}
}

func TestClonePreservesNilAndCopiesNonNil(t *testing.T) {
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
	got := Unique([]os.Signal{testSIGINT, testSIGTERM, testSIGINT, testSIGHUP})
	want := []os.Signal{testSIGINT, testSIGTERM, testSIGHUP}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if !sameSignal(got[i], want[i]) {
			t.Fatalf("signal[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestMergePreservesSetOrderAndUniqueness(t *testing.T) {
	got := Merge([]os.Signal{testSIGINT, testSIGTERM}, []os.Signal{testSIGTERM, testSIGHUP})
	want := []os.Signal{testSIGINT, testSIGTERM, testSIGHUP}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if !sameSignal(got[i], want[i]) {
			t.Fatalf("signal[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestContainsReportsSignalMembership(t *testing.T) {
	set := []os.Signal{testSIGINT, testSIGTERM}
	if !Contains(set, testSIGTERM) {
		t.Fatal("Contains did not find present signal")
	}
	if Contains(set, testSIGHUP) {
		t.Fatal("Contains found missing signal")
	}
}

func TestSignalSetHelpersRejectNilSignals(t *testing.T) {
	mustPanicWith(t, errNilSignalSetSignal, func() {
		Unique([]os.Signal{nil})
	})
	mustPanicWith(t, errNilSignalSetSignal, func() {
		Merge([]os.Signal{testSIGINT}, []os.Signal{nil})
	})
	mustPanicWith(t, errNilSignalSetSignal, func() {
		Contains([]os.Signal{testSIGINT}, nil)
	})
	mustPanicWith(t, errNilSignalSetSignal, func() {
		Contains([]os.Signal{nil}, testSIGINT)
	})
}

func TestOptionalSignalSetsAreCallable(t *testing.T) {
	_ = Reload()
	_ = Diagnostics()
}

func TestTypeNameHandlesNil(t *testing.T) {
	if got := typeName(nil); got != "<nil>" {
		t.Fatalf("typeName(nil) = %q, want <nil>", got)
	}
}
