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

// Member is one concrete object member in declaration order.
//
// Member is payload data, not a descriptor member. Name is the actual member name
// present in the object, and Value is the actual nested payload value. Empty
// names and invalid zero Values are rejected by ObjectValue.
type Member struct {
	// Name is the concrete object member name.
	Name string
	// Value is the concrete member payload.
	Value Value
}

// ObjectMember constructs one object member and clones the supplied value.
//
// Cloning here makes Member safe to pass into ObjectValue without retaining a
// mutable alias to bytes or nested composite values.
func ObjectMember(name string, value Value) Member {
	return Member{Name: name, Value: value.Clone()}
}
