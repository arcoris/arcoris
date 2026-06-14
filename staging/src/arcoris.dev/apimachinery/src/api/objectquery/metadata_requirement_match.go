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

import "slices"

// match evaluates one requirement against a key/value lookup function.
//
// Negative operators intentionally match absent keys: NotEquals is true when a
// key is missing or differs, and NotIn is true when a key is missing or outside
// the finite value set.
func (r metadataRequirement) match(lookup func(string) (string, bool)) bool {
	actual, ok := lookup(r.key)
	switch r.op {
	case OperatorExists:
		return ok
	case OperatorDoesNotExist:
		return !ok
	case OperatorEquals:
		return ok && actual == r.values[0]
	case OperatorNotEquals:
		return !ok || actual != r.values[0]
	case OperatorIn:
		return ok && slices.Contains(r.values, actual)
	case OperatorNotIn:
		return !ok || !slices.Contains(r.values, actual)
	default:
		return false
	}
}
