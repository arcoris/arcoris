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

package stamp

import "time"

// Timestamp wraps a metadata time value.
//
// It is not a causal clock and does not imply distributed ordering.
type Timestamp struct {
	// Time stores the wrapped wall-clock value with monotonic data stripped.
	Time time.Time
}

// NewTimestamp wraps t as metadata timestamp.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t.Round(0)}
}

// IsZero reports whether the timestamp is absent.
func (t Timestamp) IsZero() bool {
	return t.Time.IsZero()
}

// Clone returns a value copy of the timestamp.
func (t Timestamp) Clone() Timestamp {
	return Timestamp{Time: t.Time.Round(0)}
}
