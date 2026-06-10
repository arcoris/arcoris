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

// CanonicalText returns the canonical text form of e.
func (e SelectorEntry) CanonicalText() string {
	return strconv.Quote(e.field.String()) + ":" + e.value.CanonicalText()
}

// String returns diagnostic text for e.
func (e SelectorEntry) String() string {
	return e.CanonicalText()
}

// CanonicalText returns the canonical text form of s.
func (s Selector) CanonicalText() string {
	if len(s.entries) == 0 {
		return "{}"
	}

	parts := make([]string, len(s.entries))
	for i, entry := range s.entries {
		parts[i] = entry.CanonicalText()
	}

	return "{" + strings.Join(parts, ",") + "}"
}

// String returns diagnostic text for s.
func (s Selector) String() string {
	return s.CanonicalText()
}
