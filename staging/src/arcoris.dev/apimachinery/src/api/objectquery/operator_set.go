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

// OperatorSet is an immutable bitset of supported operators.
type OperatorSet uint64

// Operators constructs an operator support set. Unknown operators are ignored
// so callers can assemble sets defensively.
func Operators(operators ...Operator) OperatorSet {
	var set OperatorSet
	for _, op := range operators {
		if op.IsValid() {
			set |= 1 << op
		}
	}

	return set
}

// Supports reports whether op is present in set.
func (set OperatorSet) Supports(op Operator) bool {
	return op.IsValid() && set&(1<<op) != 0
}

// Built-in operator sets shared by query term domains.
var (
	// metadataOperators is the complete operator set shared by label and
	// annotation terms.
	metadataOperators = Operators(
		OperatorExists,
		OperatorDoesNotExist,
		OperatorEquals,
		OperatorNotEquals,
		OperatorIn,
		OperatorNotIn,
	)
)
