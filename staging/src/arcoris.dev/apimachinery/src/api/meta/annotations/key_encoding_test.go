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

package annotations

import (
	"encoding/json"
	"testing"
)

func TestKeyEncoding(t *testing.T) {
	text, err := Key("note").MarshalText()
	requireNoError(t, err)
	if string(text) != "note" {
		t.Fatalf("MarshalText() = %q", text)
	}

	var key Key
	requireNoError(t, key.UnmarshalText([]byte("note")))
	if key != "note" {
		t.Fatalf("UnmarshalText() = %q", key)
	}

	data, err := json.Marshal(Key("note"))
	requireNoError(t, err)
	if string(data) != `"note"` {
		t.Fatalf("MarshalJSON() = %s", data)
	}

	requireErrorIs(t, json.Unmarshal([]byte(`null`), &key), ErrInvalidJSON)
}
