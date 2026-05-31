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

import "time"

// Timestamp constructs a timestamp value and strips monotonic clock data.
//
// The location and instant are preserved, but Go's monotonic clock component is
// removed with Round(0) so equality and future serialization are deterministic.
// The zero time is allowed as concrete payload data.
func Timestamp(v time.Time) Value {
	return Value{kind: KindTimestamp, timestampValue: v.Round(0)}
}
