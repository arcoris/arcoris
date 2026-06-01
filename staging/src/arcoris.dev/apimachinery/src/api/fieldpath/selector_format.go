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
func (e SelectorEntry) String() string {
	return strconv.Quote(e.field) + ":" + e.value.String()
}

// String returns the canonical diagnostic form of s.
func (s Selector) String() string {
	if len(s.entries) == 0 {
		return "{}"
	}

	parts := make([]string, len(s.entries))
	for i, entry := range s.entries {
		parts[i] = entry.String()
	}

	return "{" + strings.Join(parts, ",") + "}"
}
