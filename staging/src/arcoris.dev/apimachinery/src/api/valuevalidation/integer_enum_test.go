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

package valuevalidation

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestValidateSignedEnumReportsMissingValue(t *testing.T) {
	run := newValidator(Options{})

	run.validateSignedEnum(fieldpath.RootPath(), 3, []int64{1, 2})

	requireInternalError(
		t,
		run.result(),
		ErrEnumMismatch,
		ErrorReasonEnumMismatch,
		"$",
	)
}

func TestValidateUnsignedEnumReportsMissingValue(t *testing.T) {
	run := newValidator(Options{})

	run.validateUnsignedEnum(fieldpath.RootPath(), 3, []uint64{1, 2})

	requireInternalError(
		t,
		run.result(),
		ErrEnumMismatch,
		ErrorReasonEnumMismatch,
		"$",
	)
}

func TestIntegerEnumAdapters(t *testing.T) {
	if got, want := signedEnum([]int8{-1, 2}), []int64{-1, 2}; !equalInt64Slice(got, want) {
		t.Fatalf("signedEnum() = %#v, want %#v", got, want)
	}
	if got, want := unsignedEnum([]uint8{1, 2}), []uint64{1, 2}; !equalUint64Slice(got, want) {
		t.Fatalf("unsignedEnum() = %#v, want %#v", got, want)
	}
}

func TestIntegerEnumContains(t *testing.T) {
	if !containsInt64([]int64{-1, 2}, -1) {
		t.Fatalf("containsInt64() = false")
	}
	if containsInt64([]int64{-1, 2}, 3) {
		t.Fatalf("containsInt64() = true")
	}
	if !containsUint64([]uint64{1, 2}, 1) {
		t.Fatalf("containsUint64() = false")
	}
	if containsUint64([]uint64{1, 2}, 3) {
		t.Fatalf("containsUint64() = true")
	}
}

func equalInt64Slice(left, right []int64) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

func equalUint64Slice(left, right []uint64) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}
