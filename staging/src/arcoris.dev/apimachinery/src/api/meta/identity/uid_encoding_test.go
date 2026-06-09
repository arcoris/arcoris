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

package identity

import (
	"encoding/json"
	"testing"
)

func TestUIDEncoding(t *testing.T) {
	text, err := UID("uid-1").MarshalText()
	requireNoError(t, err)
	if string(text) != "uid-1" {
		t.Fatalf("MarshalText() = %q", text)
	}

	var uid UID
	requireNoError(t, uid.UnmarshalText([]byte("uid-1")))
	if uid != "uid-1" {
		t.Fatalf("UnmarshalText() = %q", uid)
	}

	data, err := json.Marshal(UID("uid-1"))
	requireNoError(t, err)
	if string(data) != `"uid-1"` {
		t.Fatalf("MarshalJSON() = %s", data)
	}

	_, err = UID("").MarshalText()
	requireErrorIs(t, err, ErrInvalidUID)

	_, err = json.Marshal(UID(""))
	requireErrorIs(t, err, ErrInvalidUID)

	requireErrorIs(t, json.Unmarshal([]byte(`null`), &uid), ErrInvalidJSON)
	requireErrorIs(t, json.Unmarshal([]byte(`1`), &uid), ErrInvalidJSON)
	requireErrorIs(t, json.Unmarshal([]byte(`{}`), &uid), ErrInvalidJSON)
}
