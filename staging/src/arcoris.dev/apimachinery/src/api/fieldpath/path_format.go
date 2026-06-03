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

package fieldpath

import (
	"strconv"
	"strings"
)

// String returns the canonical diagnostic form of e.
func (e Element) String() string {
	switch e.kind {
	case ElementField:
		if isSimpleFieldName(e.name) {
			return "." + e.name
		}

		return "." + strconv.Quote(e.name)
	case ElementKey:
		return "[" + strconv.Quote(e.name) + "]"
	case ElementIndex:
		return "[" + strconv.Itoa(e.index) + "]"
	case ElementSelector:
		return "[" + e.selector.String() + "]"
	default:
		return ".<invalid>"
	}
}

// String returns the canonical diagnostic form of p.
func (p Path) String() string {
	if len(p.elements) == 0 {
		return "$"
	}

	var builder strings.Builder
	builder.Grow(1 + len(p.elements)*8)
	builder.WriteByte('$')

	for _, e := range p.elements {
		builder.WriteString(e.String())
	}

	return builder.String()
}

// isSimpleFieldName reports whether name can use dot notation.
func isSimpleFieldName(name string) bool {
	if name == "" {
		return false
	}

	for i := 0; i < len(name); i++ {
		ch := name[i]

		if i == 0 {
			if !isSimpleFieldStart(ch) {
				return false
			}
			continue
		}

		if !isSimpleFieldContinue(ch) {
			return false
		}
	}

	return true
}

// isSimpleFieldStart reports whether ch can start a dot-form field name.
func isSimpleFieldStart(ch byte) bool {
	return ch == '_' ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= 'a' && ch <= 'z')
}

// isSimpleFieldContinue reports whether ch can continue a dot-form field name.
func isSimpleFieldContinue(ch byte) bool {
	return isSimpleFieldStart(ch) || (ch >= '0' && ch <= '9')
}
