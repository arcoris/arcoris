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

import "strings"

// Equal reports whether l and other are the same selector literal.
func (l Literal) Equal(other Literal) bool {
	return l.Compare(other) == 0
}

// Compare imposes deterministic total ordering on selector literals.
func (l Literal) Compare(other Literal) int {
	switch {
	case l.kind < other.kind:
		return -1
	case l.kind > other.kind:
		return 1
	}

	switch l.kind {
	case LiteralBool:
		switch {
		case !l.boolValue && other.boolValue:
			return -1
		case l.boolValue && !other.boolValue:
			return 1
		default:
			return 0
		}
	case LiteralInteger:
		return l.intValue.compare(other.intValue)
	case LiteralString:
		return strings.Compare(l.stringValue, other.stringValue)
	default:
		return 0
	}
}
