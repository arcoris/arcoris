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
	"strings"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// jsonPath stores a syntactic JSON diagnostic path.
//
// It is deliberately not api/fieldpath.Path. JSON paths describe encoded
// document locations, not semantic API field ownership paths.
type jsonPath struct {
	// text stores the already-rendered JSON diagnostic path.
	//
	// The path helper is append-only and value-like; callers derive children by
	// returning a new jsonPath instead of mutating shared state.
	text string
}

// rootPath returns the JSON document root path.
func rootPath() jsonPath {
	return jsonPath{text: "$"}
}

// String returns the stable path text.
func (p jsonPath) String() string {
	if p.text == "" {
		return "$"
	}

	return p.text
}

// Member returns a child path for an object member.
//
// Simple identifier-like names use dot notation for readability. All other
// names use quoted bracket notation so keys containing dots, slashes, quotes,
// spaces, or ownership-style path text remain unambiguous.
func (p jsonPath) Member(name string) jsonPath {
	if isSimpleJSONPathName(name) {
		return jsonPath{text: p.String() + "." + name}
	}

	return jsonPath{text: p.String() + "[" + strconv.Quote(name) + "]"}
}

// Index returns a child path for an array item.
func (p jsonPath) Index(index int) jsonPath {
	return jsonPath{text: p.String() + "[" + strconv.Itoa(index) + "]"}
}

// isSimpleJSONPathName reports whether name can use dot notation.
//
// The rule is intentionally smaller than JSON object-name expressiveness. The
// fallback bracket form is always correct, so dot notation is reserved for
// plain identifier-looking names that are easy to scan in diagnostics.
func isSimpleJSONPathName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		if r > 0x7f {
			return false
		}
		b := byte(r)
		if i == 0 {
			if !lexical.IsASCIIAlpha(b) && b != '_' {
				return false
			}
			continue
		}
		if !lexical.IsASCIIAlpha(b) && !lexical.IsASCIIDigit(b) && b != '_' {
			return false
		}
	}

	return !strings.Contains(name, ".")
}
