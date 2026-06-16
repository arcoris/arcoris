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

package objectquery

import "strings"

// canonicalExprKey encodes expression structure into a deterministic internal
// sort key. It is deliberately not a public textual query syntax.
func canonicalExprKey(e *expr) string {
	if e == nil {
		return "all"
	}
	switch e.kind {
	case exprAll:
		return "all"
	case exprNone:
		return "none"
	case exprTerm:
		return "term(" + e.term.canonicalKey() + ")"
	case exprNot:
		return "not(" + e.children[0].canonicalKey() + ")"
	case exprAnd:
		return "and(" + joinExprKeys(e.children) + ")"
	case exprOr:
		return "or(" + joinExprKeys(e.children) + ")"
	default:
		return "unknown"
	}
}

// joinExprKeys joins child canonical keys after the caller has already
// normalized child ordering.
func joinExprKeys(children []*expr) string {
	keys := make([]string, 0, len(children))
	for _, child := range children {
		keys = append(keys, child.canonicalKey())
	}

	return strings.Join(keys, ",")
}
