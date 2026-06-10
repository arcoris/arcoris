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

package fieldpath

import "strconv"

// ValidateStructure checks whether p contains only valid semantic path elements.
//
// It does not check the path against a descriptor, object schema, ownership
// policy, apply strategy, storage layout, or wire format.
func (p Path) ValidateStructure() error {
	for i, element := range p.elements {
		if err := element.ValidateStructure(); err != nil {
			return nested(
				ErrInvalidPath,
				ErrorReasonInvalidElement,
				"path element "+strconv.Itoa(i)+" is invalid",
				err,
			)
		}
	}

	return nil
}
