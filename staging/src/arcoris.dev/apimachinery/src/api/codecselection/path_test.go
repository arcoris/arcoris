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

package codecselection

import "testing"

func TestBindingPathHelpers(t *testing.T) {
	base := "codecselection.decodeBindings[2]"

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "decode binding", got: decodeBindingPath(2), want: base},
		{name: "encode binding", got: encodeBindingPath(3), want: "codecselection.encodeBindings[3]"},
		{name: "content type", got: bindingContentTypePath(base), want: base + ".contentType"},
		{name: "target", got: bindingTargetPath(base), want: base + ".target"},
		{name: "transport", got: bindingTransportPath(base), want: base + ".transport"},
		{name: "entry ID", got: bindingEntryIDPath(base), want: base + ".entryID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("path = %q; want %q", tt.got, tt.want)
			}
		})
	}
}

func TestRuntimePathHelpers(t *testing.T) {
	if got := runtimeDecodePath("object"); got != "codecselection.decode.object" {
		t.Fatalf("runtimeDecodePath() = %q; want codecselection.decode.object", got)
	}
	if got := runtimeEncodePath("object"); got != "codecselection.encode.object" {
		t.Fatalf("runtimeEncodePath() = %q; want codecselection.encode.object", got)
	}
}
