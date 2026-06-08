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

import "testing"

func TestErrorReasonValues(t *testing.T) {
	tests := map[ErrorReason]string{
		ErrorReasonInvalidContentType:      "invalid_content_type",
		ErrorReasonInvalidParameters:       "invalid_parameters",
		ErrorReasonInvalidPreference:       "invalid_preference",
		ErrorReasonInvalidBinding:          "invalid_binding",
		ErrorReasonDuplicateDecodeBinding:  "duplicate_decode_binding",
		ErrorReasonDuplicateEncodeBinding:  "duplicate_encode_binding",
		ErrorReasonUnknownEntryID:          "unknown_entry_id",
		ErrorReasonEntryMediaTypeMismatch:  "entry_media_type_mismatch",
		ErrorReasonEntryTargetMismatch:     "entry_target_mismatch",
		ErrorReasonEntryCapabilityMismatch: "entry_capability_mismatch",
		ErrorReasonNoDecodeBinding:         "no_decode_binding",
		ErrorReasonNoEncodePreference:      "no_encode_preference",
	}

	for reason, want := range tests {
		t.Run(string(reason), func(t *testing.T) {
			if string(reason) != want {
				t.Fatalf("reason = %q; want %q", reason, want)
			}
		})
	}
}
