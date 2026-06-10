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
)

// CanonicalText returns the canonical text form of l.
func (l Literal) CanonicalText() string {
	switch l.kind {
	case LiteralBool:
		return strconv.FormatBool(l.boolValue)
	case LiteralInteger:
		return l.intValue.string()
	case LiteralString:
		return strconv.Quote(l.stringValue)
	default:
		return "<invalid>"
	}
}

// String returns diagnostic text for l.
func (l Literal) String() string {
	return l.CanonicalText()
}
