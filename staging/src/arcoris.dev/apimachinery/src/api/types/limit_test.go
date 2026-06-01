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

package types

import "testing"

func TestInvalidRange(t *testing.T) {
	requireEqual(t, invalidRange(limit[int8]{}, limit[int8]{}), false)
	requireEqual(t, invalidRange(limit[int8]{value: 1, set: true}, limit[int8]{}), false)
	requireEqual(t, invalidRange(limit[uint64]{}, limit[uint64]{value: 1, set: true}), false)
	requireEqual(
		t,
		invalidRange(limit[int8]{value: 1, set: true}, limit[int8]{value: 1, set: true}),
		false,
	)
	requireEqual(
		t,
		invalidRange(limit[int8]{value: 2, set: true}, limit[int8]{value: 1, set: true}),
		true,
	)
	requireEqual(
		t,
		invalidRange(limit[uint64]{value: 2, set: true}, limit[uint64]{value: 1, set: true}),
		true,
	)
	requireEqual(
		t,
		invalidRange(limit[string]{value: "b", set: true}, limit[string]{value: "a", set: true}),
		true,
	)
}

func TestValidateLengthLimits(t *testing.T) {
	requireNoError(
		t,
		validateLengthLimits(
			limit[int]{value: 0, set: true},
			limit[int]{value: 1, set: true},
			"len",
		),
	)
	requireErrorIs(
		t,
		validateLengthLimits(limit[int]{value: -1, set: true}, limit[int]{}, "len"),
		ErrInvalidType,
	)
	requireErrorIs(
		t,
		validateLengthLimits(limit[int]{}, limit[int]{value: -1, set: true}, "len"),
		ErrInvalidType,
	)
	requireErrorIs(
		t,
		validateLengthLimits(
			limit[int]{value: 2, set: true},
			limit[int]{value: 1, set: true},
			"len",
		),
		ErrInvalidType,
	)
}
