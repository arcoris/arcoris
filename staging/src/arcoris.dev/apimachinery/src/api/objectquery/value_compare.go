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

import "arcoris.dev/apimachinery/api/value"

// valueEqual delegates concrete payload equality to api/value so objectquery
// does not duplicate value semantics.
func valueEqual(left value.Value, right value.Value) bool {
	return value.Equal(left, right)
}

// valueIn applies equality semantics to a canonical literal set.
func valueIn(actual value.Value, literals []value.Value) bool {
	actualKey := canonicalValueKey(actual)
	for _, literal := range literals {
		if canonicalValueKey(literal) == actualKey {
			return true
		}
	}

	return false
}
