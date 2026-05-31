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

package value

// TimeOfDayValue constructs a time-of-day Value after checking range invariants.
//
// TimeOfDay is immutable by convention and contains no mutable backing data, so
// the payload can be copied directly after validity is confirmed. The zero
// TimeOfDay is valid and represents midnight.
func TimeOfDayValue(v TimeOfDay) (Value, error) {
	if !v.IsValid() {
		return Value{}, invalidTimeOfDay("time-of-day value is invalid")
	}

	return Value{kind: KindTimeOfDay, timeOfDayValue: v}, nil
}

// MustTimeOfDayValue constructs a time-of-day Value or panics when v is invalid.
//
// It is intended for tests and static fixtures where an invalid TimeOfDay is a
// programmer error. Runtime construction paths should use TimeOfDayValue.
func MustTimeOfDayValue(v TimeOfDay) Value {
	value, err := TimeOfDayValue(v)
	if err != nil {
		panic(err)
	}

	return value
}
