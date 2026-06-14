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

// ErrorReason identifies a precise object query validation failure.
type ErrorReason string

const (
	// ErrorReasonInvalidQuery identifies a malformed top-level query.
	ErrorReasonInvalidQuery ErrorReason = "invalid_query"

	// ErrorReasonInvalidIdentity identifies malformed identity requirements.
	ErrorReasonInvalidIdentity ErrorReason = "invalid_identity"

	// ErrorReasonInvalidSelector identifies malformed metadata selector state.
	ErrorReasonInvalidSelector ErrorReason = "invalid_selector"

	// ErrorReasonInvalidRequirement identifies malformed metadata requirements.
	ErrorReasonInvalidRequirement ErrorReason = "invalid_requirement"

	// ErrorReasonInvalidOperator identifies an unknown operator.
	ErrorReasonInvalidOperator ErrorReason = "invalid_operator"

	// ErrorReasonInvalidKey identifies an invalid metadata key.
	ErrorReasonInvalidKey ErrorReason = "invalid_key"

	// ErrorReasonInvalidValue identifies an invalid metadata value.
	ErrorReasonInvalidValue ErrorReason = "invalid_value"

	// ErrorReasonInvalidValueCount identifies an operator/value arity mismatch.
	ErrorReasonInvalidValueCount ErrorReason = "invalid_value_count"
)
