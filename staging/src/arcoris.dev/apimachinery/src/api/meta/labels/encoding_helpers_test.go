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

import "testing"

func TestEncodingHelpers(t *testing.T) {
	text, err := marshalText("role", func() error { return nil })
	requireNoError(t, err)
	if string(text) != "role" {
		t.Fatalf("marshalText() = %q", text)
	}

	data, err := marshalJSONString("role", func() error { return nil })
	requireNoError(t, err)
	if string(data) != `"role"` {
		t.Fatalf("marshalJSONString() = %s", data)
	}

	value, err := unmarshalJSONString("label.key", []byte(`"role"`))
	requireNoError(t, err)
	if value != "role" {
		t.Fatalf("unmarshalJSONString() = %q", value)
	}

	requireErrorIs(t, nilReceiver("label.key"), ErrNilReceiver)
	requireErrorIs(t, marshalError(), ErrInvalidKey)
	requireErrorIs(t, unmarshalJSONError(), ErrInvalidJSON)
}

func marshalError() error {
	_, err := marshalText("bad", func() error {
		return invalid("label.key", ErrInvalidKey, ErrorReasonInvalidForm, "bad key")
	})
	return err
}

func unmarshalJSONError() error {
	_, err := unmarshalJSONString("label.key", []byte(`null`))
	return err
}
