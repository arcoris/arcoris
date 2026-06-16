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

import "testing"

// TestErrorHelpersPreserveSentinels verifies representative error builders
// retain both broad and specific errors.Is classifications.
func TestErrorHelpersPreserveSentinels(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
	}{
		{name: "expression", err: invalidExpressionError("bad"), target: ErrInvalidExpression},
		{name: "term", err: invalidTermError("bad"), target: ErrInvalidTerm},
		{name: "field", err: invalidFieldError(ErrInvalidTerm, "bad"), target: ErrInvalidField},
		{name: "operator", err: invalidOperatorError(Operator(99)), target: ErrInvalidOperator},
		{name: "unsupported", err: unsupportedOperatorError(OperatorContains, "metadata"), target: ErrUnsupportedOperator},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, ErrInvalidQuery)
			requireErrorIs(t, tt.err, tt.target)
		})
	}
}
