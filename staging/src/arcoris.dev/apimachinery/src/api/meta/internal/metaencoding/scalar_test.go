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

package metaencoding

import (
	"errors"
	"testing"
)

func TestMarshalText(t *testing.T) {
	data, err := MarshalText("value", func() error { return nil })
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}
	if string(data) != "value" {
		t.Fatalf("MarshalText() = %q", data)
	}

	want := errors.New("invalid")
	if _, err := MarshalText("value", func() error { return want }); !errors.Is(err, want) {
		t.Fatalf("MarshalText() error = %v, want %v", err, want)
	}
}

func TestMarshalJSONString(t *testing.T) {
	data, err := MarshalJSONString("value", func() error { return nil })
	if err != nil {
		t.Fatalf("MarshalJSONString() error = %v", err)
	}
	if string(data) != `"value"` {
		t.Fatalf("MarshalJSONString() = %s", data)
	}

	want := errors.New("invalid")
	if _, err := MarshalJSONString("value", func() error { return want }); !errors.Is(err, want) {
		t.Fatalf("MarshalJSONString() error = %v, want %v", err, want)
	}
}

func TestDecodeJSONString(t *testing.T) {
	value, isNull, err := DecodeJSONString([]byte(`"value"`))
	if err != nil {
		t.Fatalf("DecodeJSONString() error = %v", err)
	}
	if isNull || value != "value" {
		t.Fatalf("DecodeJSONString() = %q, %v", value, isNull)
	}

	_, isNull, err = DecodeJSONString([]byte(`null`))
	if err != nil {
		t.Fatalf("DecodeJSONString(null) error = %v", err)
	}
	if !isNull {
		t.Fatal("DecodeJSONString(null) isNull = false")
	}

	if _, _, err := DecodeJSONString([]byte(`{}`)); err == nil {
		t.Fatal("DecodeJSONString(object) error = nil")
	}
}
