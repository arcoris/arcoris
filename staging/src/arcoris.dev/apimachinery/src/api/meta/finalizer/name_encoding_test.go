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

package finalizer

import (
	"encoding/json"
	"testing"
)

func TestNameEncoding(t *testing.T) {
	text, err := Name("cleanup").MarshalText()
	requireNoError(t, err)
	if string(text) != "cleanup" {
		t.Fatalf("MarshalText() = %q", text)
	}

	var name Name
	requireNoError(t, name.UnmarshalText([]byte("cleanup")))
	if name != "cleanup" {
		t.Fatalf("UnmarshalText() = %q", name)
	}

	data, err := json.Marshal(Name("cleanup"))
	requireNoError(t, err)
	if string(data) != `"cleanup"` {
		t.Fatalf("MarshalJSON() = %s", data)
	}

	requireErrorIs(t, json.Unmarshal([]byte(`null`), &name), ErrInvalidJSON)
}
