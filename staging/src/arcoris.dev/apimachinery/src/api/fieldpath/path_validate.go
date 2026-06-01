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

// Validate checks whether e is a valid semantic path element.
func (e Element) Validate() error {
	switch e.kind {
	case ElementField:
		if e.name == "" {
			return nested(
				ErrInvalidElement,
				ErrorReasonEmptyFieldName,
				"field element name is empty",
				ErrEmptyFieldName,
			)
		}
	case ElementKey:
		if e.name == "" {
			return nested(
				ErrInvalidElement,
				ErrorReasonEmptyKey,
				"key element is empty",
				ErrEmptyKey,
			)
		}
	case ElementIndex:
		if e.index < 0 {
			return nested(
				ErrInvalidElement,
				ErrorReasonNegativeIndex,
				"index element is negative",
				ErrNegativeIndex,
			)
		}
	case ElementSelector:
		if err := e.selector.Validate(); err != nil {
			return nested(
				ErrInvalidElement,
				ErrorReasonInvalidSelector,
				"selector element is invalid",
				err,
			)
		}
	default:
		return errorf(
			ErrInvalidElement,
			ErrorReasonInvalidElement,
			"element kind is invalid",
		)
	}

	return nil
}

// Validate checks whether p contains only valid semantic path elements.
func (p Path) Validate() error {
	for i, element := range p.elements {
		if err := element.Validate(); err != nil {
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
