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

package codecjson

import (
	"strconv"
	"testing"
)

func TestJSONPathRoot(t *testing.T) {
	if got := rootPath().String(); got != "$" {
		t.Fatalf("root path = %q; want $", got)
	}
}

func TestJSONPathSimpleMember(t *testing.T) {
	if got := rootPath().Member("desired").String(); got != "$.desired" {
		t.Fatalf("member path = %q", got)
	}
}

func TestJSONPathQuotedMember(t *testing.T) {
	got := rootPath().Member("example.com/key").String()
	if got != `$["example.com/key"]` {
		t.Fatalf("member path = %q", got)
	}
}

func TestJSONPathQuotedNonASCIIMember(t *testing.T) {
	name := string([]rune{0x043a, 0x043b, 0x044e, 0x0447})
	want := "$[" + strconv.Quote(name) + "]"

	got := rootPath().Member(name).String()
	if got != want {
		t.Fatalf("member path = %q", got)
	}
}

func TestJSONPathIndex(t *testing.T) {
	if got := rootPath().Member("items").Index(3).String(); got != "$.items[3]" {
		t.Fatalf("index path = %q", got)
	}
}
