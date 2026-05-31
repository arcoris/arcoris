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

// Entry is one concrete map entry in caller-supplied order.
//
// Entry is payload data, not a descriptor. Empty keys and invalid zero Values
// are rejected by NewMap.
type Entry struct {
	// Key is the concrete map entry key.
	Key string
	// Value is the concrete map entry payload.
	Value Value
}

// MapEntry constructs one map entry and clones the supplied value.
//
// Cloning here makes Entry safe to pass into NewMap without retaining a mutable
// alias to bytes or nested composite values.
func MapEntry(key string, value Value) Entry {
	return Entry{Key: key, Value: value.Clone()}
}
