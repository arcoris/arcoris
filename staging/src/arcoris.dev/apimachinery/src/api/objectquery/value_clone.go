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

// cloneValues detaches a literal slice and every contained value.
func cloneValues(values []value.Value) []value.Value {
	if len(values) == 0 {
		return nil
	}

	out := make([]value.Value, len(values))
	for i, val := range values {
		out[i] = val.Clone()
	}

	return out
}
