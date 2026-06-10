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

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	tests := []struct {
		reason ErrorReason
		want   string
	}{
		{reason: ErrorReasonInvalidValue, want: "invalid_value"},
		{reason: ErrorReasonInvalidRecordMember, want: "invalid_record_member"},
		{reason: ErrorReasonDuplicateMemberName, want: "duplicate_member_name"},
		{reason: ErrorReasonEmptyMemberName, want: "empty_member_name"},
		{reason: ErrorReasonInvalidFloat, want: "invalid_float"},
		{reason: ErrorReasonInvalidDecimal, want: "invalid_decimal"},
		{reason: ErrorReasonInvalidDate, want: "invalid_date"},
		{reason: ErrorReasonInvalidTime, want: "invalid_time"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			requireEqual(t, string(tt.reason), tt.want)
		})
	}
}
