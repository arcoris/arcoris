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

package labels

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestSetJSONEncoding(t *testing.T) {
	data, err := json.Marshal(Set{"app": "scheduler"})
	requireNoError(t, err)
	if string(data) != `{"app":"scheduler"}` {
		t.Fatalf("json = %s", data)
	}

	var set Set
	requireNoError(t, json.Unmarshal([]byte(`{"app":"scheduler"}`), &set))
	if value, ok := set.Get("app"); !ok || value != "scheduler" {
		t.Fatalf("unmarshaled set = %#v", set)
	}
}

func TestSetJSONRejectsInvalidMetadata(t *testing.T) {
	_, err := json.Marshal(Set{"App": "scheduler"})
	requireErrorIs(t, err, ErrInvalidSet)

	var set Set
	err = json.Unmarshal([]byte(`{"App":"scheduler"}`), &set)
	requireErrorIs(t, err, ErrInvalidSet)

	err = json.Unmarshal([]byte(`{"app":1}`), &set)
	requireErrorIs(t, err, ErrInvalidJSON)
}

func TestSetJSONNilReceiver(t *testing.T) {
	var set *Set
	err := set.UnmarshalJSON([]byte(`{"app":"scheduler"}`))
	requireErrorIs(t, err, ErrNilReceiver)

	var labelErr *Error
	if !errors.As(err, &labelErr) {
		t.Fatalf("errors.As(%T) = false", labelErr)
	}
	if labelErr.Path != "labels" {
		t.Fatalf("Path = %q", labelErr.Path)
	}
}
