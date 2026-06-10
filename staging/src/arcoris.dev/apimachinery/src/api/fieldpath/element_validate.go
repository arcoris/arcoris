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

// ValidateStructure checks whether e is a well-formed semantic path element.
//
// Element validation is intentionally local to the element itself. Path
// validation only walks the element sequence and wraps the failing position,
// while selector validation remains owned by Selector.
func (e Element) ValidateStructure() error {
	switch e.kind {
	case ElementField:
		return validateFieldElement(e)
	case ElementKey:
		return validateKeyElement(e)
	case ElementIndex:
		return validateIndexElement(e)
	case ElementSelector:
		return validateSelectorElement(e)
	default:
		return newError(
			ErrInvalidElement,
			ErrorReasonInvalidElement,
			"element kind is invalid",
		)
	}
}

// validateFieldElement enforces the base grammar for fixed object-field steps.
func validateFieldElement(e Element) error {
	return e.field.ValidateStructure()
}

// validateKeyElement enforces the base grammar for dynamic map-key steps.
func validateKeyElement(e Element) error {
	return e.key.ValidateStructure()
}

// validateIndexElement rejects negative physical list positions.
func validateIndexElement(e Element) error {
	if e.index >= 0 {
		return nil
	}

	return nested(
		ErrInvalidElement,
		ErrorReasonNegativeIndex,
		"index element is negative",
		ErrNegativeIndex,
	)
}

// validateSelectorElement delegates associative-list identity validation to
// Selector while preserving the element-level error class.
func validateSelectorElement(e Element) error {
	if err := e.selector.ValidateStructure(); err != nil {
		return nested(
			ErrInvalidElement,
			ErrorReasonInvalidSelector,
			"selector element is invalid",
			err,
		)
	}

	return nil
}
