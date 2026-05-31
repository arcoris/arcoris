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

package value

import "testing"

func TestIntegerValueConstructors(t *testing.T) {
	tests := []Value{
		Int64(-1),
		Uint64(1),
		IntegerValue(Integer{negative: true}),
	}

	for _, value := range tests {
		integer, ok := value.Integer()
		requireEqual(t, ok, true)
		requireEqual(t, value.Kind(), KindInteger)
		requireEqual(t, integer.Magnitude() >= 0, true)
	}
}
