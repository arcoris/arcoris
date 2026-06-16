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

import (
	"fmt"
	"sort"

	"arcoris.dev/apimachinery/api/value"
)

// canonicalValues clones, sorts, and deduplicates literal sets by semantic
// value equality. Zero values are rejected because they do not represent an
// intentional query literal.
func canonicalValues(values []value.Value) ([]value.Value, error) {
	if len(values) == 0 {
		return nil, nil
	}

	type keyed struct {
		key   string
		value value.Value
	}
	keyedValues := make([]keyed, 0, len(values))
	for _, literal := range values {
		if literal.IsZero() {
			return nil, fmt.Errorf("%w: zero literal", ErrInvalidTerm)
		}
		keyedValues = append(keyedValues, keyed{
			key:   canonicalValueKey(literal),
			value: literal.Clone(),
		})
	}
	sort.Slice(keyedValues, func(i int, j int) bool {
		return keyedValues[i].key < keyedValues[j].key
	})

	out := make([]value.Value, 0, len(keyedValues))
	var last string
	for i, item := range keyedValues {
		if i == 0 || item.key != last {
			out = append(out, item.value)
			last = item.key
		}
	}

	return out, nil
}
