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
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestSelectableFieldValidate verifies registry entries must provide a valid
// ref, concrete value kind, and non-empty operator set.
func TestSelectableFieldValidate(t *testing.T) {
	requireNoError(t, selectable(fieldRef("spec.image"), value.KindString, Operators(OperatorEquals)).Validate())

	err := selectable(FieldRef{}, value.KindString, Operators(OperatorEquals)).Validate()
	requireErrorIs(t, err, ErrInvalidField)

	err = selectable(fieldRef("spec.image"), value.KindInvalid, Operators(OperatorEquals)).Validate()
	requireErrorIs(t, err, ErrInvalidField)

	err = selectable(fieldRef("spec.image"), value.KindString, 0).Validate()
	requireErrorIs(t, err, ErrInvalidField)
}
