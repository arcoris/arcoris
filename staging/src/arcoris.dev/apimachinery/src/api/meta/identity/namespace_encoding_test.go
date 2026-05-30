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

func TestNamespaceEncoding(t *testing.T) {
	text, err := Namespace("system").MarshalText()
	requireNoError(t, err)
	if string(text) != "system" {
		t.Fatalf("MarshalText() = %q", text)
	}

	var namespace Namespace
	requireNoError(t, namespace.UnmarshalText([]byte("system")))
	if namespace != "system" {
		t.Fatalf("UnmarshalText() = %q", namespace)
	}

	data, err := json.Marshal(Namespace("system"))
	requireNoError(t, err)
	if string(data) != `"system"` {
		t.Fatalf("MarshalJSON() = %s", data)
	}
}
