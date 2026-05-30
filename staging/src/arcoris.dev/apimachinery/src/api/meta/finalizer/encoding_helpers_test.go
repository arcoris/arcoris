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

import "testing"

func TestEncodingHelpers(t *testing.T) {
	text, err := marshalText("cleanup", func() error { return nil })
	requireNoError(t, err)
	if string(text) != "cleanup" {
		t.Fatalf("marshalText() = %q", text)
	}

	data, err := marshalJSONString("cleanup", func() error { return nil })
	requireNoError(t, err)
	if string(data) != `"cleanup"` {
		t.Fatalf("marshalJSONString() = %s", data)
	}

	value, err := unmarshalJSONString("finalizer.name", []byte(`"cleanup"`))
	requireNoError(t, err)
	if value != "cleanup" {
		t.Fatalf("unmarshalJSONString() = %q", value)
	}

	requireErrorIs(t, nilReceiver("finalizer.name"), ErrNilReceiver)
	requireErrorIs(t, marshalError(), ErrInvalidName)
	requireErrorIs(t, unmarshalJSONError(), ErrInvalidJSON)
}

func marshalError() error {
	_, err := marshalText("bad", func() error {
		return invalid("finalizer.name", ErrInvalidName, ErrorReasonInvalidForm, "bad name")
	})
	return err
}

func unmarshalJSONError() error {
	_, err := unmarshalJSONString("finalizer.name", []byte(`null`))
	return err
}
