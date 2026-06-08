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

package codecselection

import (
	"strings"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// normalizeParameterAt validates p and returns its normalized form.
func normalizeParameterAt(path string, p Parameter) (Parameter, error) {
	name := strings.ToLower(strings.TrimSpace(p.name))
	value := strings.TrimSpace(p.value)

	if violation := lexical.ValidateASCIIToken(name, parameterNameTokenOptions()); violation != nil {
		return Parameter{}, errorfAt(
			path+".name",
			ErrInvalidParameters,
			ErrorReasonInvalidParameters,
			"parameter name is invalid: %s",
			violation.Detail,
		)
	}
	if violation := lexical.ValidateASCIIToken(value, parameterValueTokenOptions()); violation != nil {
		return Parameter{}, errorfAt(
			path+".value",
			ErrInvalidParameters,
			ErrorReasonInvalidParameters,
			"parameter value is invalid: %s",
			violation.Detail,
		)
	}

	return Parameter{name: name, value: value}, nil
}
