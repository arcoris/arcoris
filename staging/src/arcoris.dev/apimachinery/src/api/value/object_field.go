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

// Field is one concrete object field in declaration order.
//
// Field is payload data, not a descriptor field. Name is the actual field name
// present in the object, and Value is the actual nested payload value. Empty
// names and invalid zero Values are rejected by NewObject.
type Field struct {
	// Name is the concrete object field name.
	Name string
	// Value is the concrete field payload.
	Value Value
}

// ObjectField constructs one object field and clones the supplied value.
//
// Cloning here makes Field safe to pass into NewObject without retaining a
// mutable alias to bytes or nested composite values.
func ObjectField(name string, value Value) Field {
	return Field{Name: name, Value: value.Clone()}
}
