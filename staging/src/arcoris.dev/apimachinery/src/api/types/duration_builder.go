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

// Duration starts a field descriptor for a duration value.
//
// The descriptor does not choose ISO-8601 duration text, Go-style duration
// text, or integer nanoseconds. Those choices belong to codecs and exporters.
//
// Field builder flow:
//
//	Field("timeout").Duration().
//		Optional().
//		Nullable().
//		Description("Operation timeout.")
func (b FieldBuilder) Duration() DurationField {
	return DurationField{field: b.state(), typ: Duration()}
}
