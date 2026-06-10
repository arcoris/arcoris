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
	"errors"
	"testing"
)

func TestSetJSONEncoding(t *testing.T) {
	data, err := json.Marshal(Set{"note": "human readable"})
	requireNoError(t, err)
	if string(data) != `{"note":"human readable"}` {
		t.Fatalf("json = %s", data)
	}

	var set Set
	requireNoError(t, json.Unmarshal([]byte(`{"note":"human readable"}`), &set))
	if value, ok := set.Get("note"); !ok || value != "human readable" {
		t.Fatalf("unmarshaled set = %#v", set)
	}
}

func TestSetJSONRejectsInvalidMetadata(t *testing.T) {
	_, err := json.Marshal(Set{"Note": "human readable"})
	requireErrorIs(t, err, ErrInvalidSet)

	var set Set
	err = json.Unmarshal([]byte(`{"Note":"human readable"}`), &set)
	requireErrorIs(t, err, ErrInvalidSet)

	err = json.Unmarshal([]byte(`{"note":1}`), &set)
	requireErrorIs(t, err, ErrInvalidJSON)
}

func TestSetJSONNilReceiver(t *testing.T) {
	var set *Set
	err := set.UnmarshalJSON([]byte(`{"note":"human readable"}`))
	requireErrorIs(t, err, ErrNilReceiver)

	var annotationErr *Error
	if !errors.As(err, &annotationErr) {
		t.Fatalf("errors.As(%T) = false", annotationErr)
	}
	if annotationErr.Path != "annotations" {
		t.Fatalf("Path = %q", annotationErr.Path)
	}
}
