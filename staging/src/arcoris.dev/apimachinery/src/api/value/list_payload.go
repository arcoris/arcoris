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

// listPayload stores ordered concrete item values.
//
// The payload stores only values. Element type, min/max length, uniqueness, and
// collection strategy are descriptor-level concerns.
type listPayload struct {
	// items contains cloned item values in stable caller order.
	items []Value
}

// newListPayload validates items and clones nested values.
//
// Invalid zero Values are rejected so every stored item is a concrete payload
// value. Explicit Null is allowed because it is real payload data.
func newListPayload(items []Value) (listPayload, error) {
	if len(items) == 0 {
		return listPayload{}, nil
	}

	payload := listPayload{
		items: make([]Value, 0, len(items)),
	}

	for i, item := range items {
		if err := validateListItem(i, item); err != nil {
			return listPayload{}, err
		}

		payload.items = append(payload.items, item.Clone())
	}

	return payload, nil
}
