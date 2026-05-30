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

package stamp

import (
	"encoding/json"
	"testing"
)

func TestResourceVersionEncoding(t *testing.T) {
	text, err := ResourceVersion("rv-1").MarshalText()
	requireNoError(t, err)
	if string(text) != "rv-1" {
		t.Fatalf("MarshalText() = %q", text)
	}

	var version ResourceVersion
	requireNoError(t, version.UnmarshalText([]byte("rv-1")))
	if version != "rv-1" {
		t.Fatalf("UnmarshalText() = %q", version)
	}

	data, err := json.Marshal(ResourceVersion("rv-1"))
	requireNoError(t, err)
	if string(data) != `"rv-1"` {
		t.Fatalf("MarshalJSON() = %s", data)
	}
}
