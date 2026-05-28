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

// Time starts a field descriptor for a time-of-day value.
//
// The descriptor records the temporal kind only. Clock precision, timezone
// interpretation, and text/binary wire formats are codec concerns.
//
// Field builder flow:
//
//	Field("startTime").Time().
//		Optional().
//		Nullable().
//		Description("Start time of day.")
func (b FieldBuilder) Time() TimeField {
	return TimeField{field: b.state(), typ: Time()}
}
