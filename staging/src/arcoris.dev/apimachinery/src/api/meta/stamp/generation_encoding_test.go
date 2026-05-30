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
)

func TestGenerationEncoding(t *testing.T) {
	text, err := Generation(12).MarshalText()
	requireNoError(t, err)
	if string(text) != "12" {
		t.Fatalf("MarshalText() = %q", text)
	}

	var generation Generation
	requireNoError(t, generation.UnmarshalText([]byte("12")))
	if generation != 12 {
		t.Fatalf("UnmarshalText() = %d", generation)
	}

	data, err := json.Marshal(Generation(12))
	requireNoError(t, err)
	if string(data) != "12" {
		t.Fatalf("MarshalJSON() = %s", data)
	}
}
