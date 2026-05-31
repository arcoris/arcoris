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

// ErrorReason identifies a precise API value construction failure.
type ErrorReason string

// Error reasons refine broad construction sentinels with stable diagnostics.
const (
	// ErrorReasonInvalidValue identifies an invalid nested Value input.
	ErrorReasonInvalidValue ErrorReason = "invalid_value"
	// ErrorReasonInvalidField identifies a malformed object field input.
	ErrorReasonInvalidField ErrorReason = "invalid_field"
	// ErrorReasonInvalidEntry identifies a malformed map entry input.
	ErrorReasonInvalidEntry ErrorReason = "invalid_entry"
	// ErrorReasonDuplicateName identifies a repeated object field name.
	ErrorReasonDuplicateName ErrorReason = "duplicate_name"
	// ErrorReasonDuplicateKey identifies a repeated map entry key.
	ErrorReasonDuplicateKey ErrorReason = "duplicate_key"
	// ErrorReasonEmptyName identifies an empty object field name.
	ErrorReasonEmptyName ErrorReason = "empty_name"
	// ErrorReasonEmptyKey identifies an empty map entry key.
	ErrorReasonEmptyKey ErrorReason = "empty_key"
	// ErrorReasonInvalidFloat identifies a non-finite float input.
	ErrorReasonInvalidFloat ErrorReason = "invalid_float"
	// ErrorReasonInvalidDecimal identifies malformed decimal text.
	ErrorReasonInvalidDecimal ErrorReason = "invalid_decimal"
	// ErrorReasonInvalidDate identifies an impossible calendar date.
	ErrorReasonInvalidDate ErrorReason = "invalid_date"
	// ErrorReasonInvalidTime identifies an impossible time-of-day.
	ErrorReasonInvalidTime ErrorReason = "invalid_time"
)
