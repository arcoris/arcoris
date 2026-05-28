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

// TimeType builds descriptors for times of day without calendar dates.
//
// TimeType records a time-of-day descriptor without a calendar date. It does
// not choose an external textual or numeric encoding.
type TimeType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact time-of-day constraints under construction.
	payload timePayload
}

// Time returns a time-of-day descriptor builder.
//
// Typical reusable declaration:
//
//	startTimeType := Time().Nullable()
func Time() TimeType {
	return TimeType{header: newHeader(TypeTime)}
}

// Nullable returns a time descriptor that admits null values.
func (t TimeType) Nullable() TimeType {
	t.header = t.header.withNullable()
	return t
}

// Type returns a detached Type descriptor.
func (t TimeType) Type() Type {
	out := typeFromHeader(t.header)
	out.timeOfDay = cloneTimePayload(t.payload)
	return out
}

// typeExpr marks TimeType as a sealed TypeExpr implementation.
func (t TimeType) typeExpr() {}
