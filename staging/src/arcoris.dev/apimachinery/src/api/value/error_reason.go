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
	// ErrorReasonInvalidRecordMember identifies a malformed record member input.
	ErrorReasonInvalidRecordMember ErrorReason = "invalid_record_member"
	// ErrorReasonDuplicateMemberName identifies a repeated record member name.
	ErrorReasonDuplicateMemberName ErrorReason = "duplicate_member_name"
	// ErrorReasonEmptyMemberName identifies an empty record member name.
	ErrorReasonEmptyMemberName ErrorReason = "empty_member_name"
	// ErrorReasonInvalidFloat identifies a non-finite float input.
	ErrorReasonInvalidFloat ErrorReason = "invalid_float"
	// ErrorReasonInvalidDecimal identifies malformed decimal text.
	ErrorReasonInvalidDecimal ErrorReason = "invalid_decimal"
	// ErrorReasonInvalidDate identifies an impossible calendar date.
	ErrorReasonInvalidDate ErrorReason = "invalid_date"
	// ErrorReasonInvalidTime identifies an impossible time-of-day.
	ErrorReasonInvalidTime ErrorReason = "invalid_time"
)
