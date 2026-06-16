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

import "arcoris.dev/apimachinery/api/value"

// constraintsForMetadataTerm emits only positive metadata constraints that can
// narrow candidates without changing absent-key semantics.
func constraintsForMetadataTerm(t term) []Constraint {
	switch t.operator {
	case OperatorExists, OperatorEquals, OperatorIn:
	default:
		return nil
	}

	kind := ConstraintLabel
	if t.metadataDomain == metadataAnnotations {
		kind = ConstraintAnnotation
	}

	values := make([]value.Value, 0, len(t.stringValues))
	for _, val := range t.stringValues {
		values = append(values, value.StringValue(val))
	}

	return []Constraint{{
		Kind:   kind,
		Ref:    ConstraintRef{Metadata: t.metadataKey},
		Op:     t.operator,
		Values: values,
	}}
}
