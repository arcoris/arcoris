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

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestampEncoding(t *testing.T) {
	timestamp := NewTimestamp(time.Date(2026, 5, 30, 12, 0, 0, 123, time.UTC))

	text, err := timestamp.MarshalText()
	requireNoError(t, err)

	var parsed Timestamp
	requireNoError(t, parsed.UnmarshalText(text))
	if parsed.IsZero() {
		t.Fatal("UnmarshalText() produced zero timestamp")
	}

	data, err := json.Marshal(timestamp)
	requireNoError(t, err)
	requireNoError(t, json.Unmarshal(data, &parsed))
}
