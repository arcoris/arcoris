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

import "testing"

func TestEncodingHelpers(t *testing.T) {
	text, err := marshalText("worker", func() error { return nil })
	requireNoError(t, err)
	if string(text) != "worker" {
		t.Fatalf("marshalText() = %q", text)
	}

	data, err := marshalJSONString("worker", func() error { return nil })
	requireNoError(t, err)
	if string(data) != `"worker"` {
		t.Fatalf("marshalJSONString() = %s", data)
	}

	value, err := unmarshalJSONString("name", []byte(`"worker"`))
	requireNoError(t, err)
	if value != "worker" {
		t.Fatalf("unmarshalJSONString() = %q", value)
	}

	requireErrorIs(t, nilReceiver("name"), ErrNilReceiver)
	requireErrorIs(t, marshalError(), ErrInvalidName)
	requireErrorIs(t, unmarshalJSONError(), ErrInvalidJSON)
}

func marshalError() error {
	_, err := marshalText("bad", func() error {
		return invalid("name", ErrInvalidName, ErrorReasonInvalidForm, "bad name")
	})
	return err
}

func unmarshalJSONError() error {
	_, err := unmarshalJSONString("name", []byte(`null`))
	return err
}
