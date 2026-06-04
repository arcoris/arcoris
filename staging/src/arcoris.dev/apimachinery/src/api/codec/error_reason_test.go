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

package codec

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	testCases := map[ErrorReason]string{
		ErrorReasonInvalidFormat:        "invalid_format",
		ErrorReasonInvalidMediaType:     "invalid_media_type",
		ErrorReasonInvalidTarget:        "invalid_target",
		ErrorReasonInvalidInfo:          "invalid_info",
		ErrorReasonUnsupportedFormat:    "unsupported_format",
		ErrorReasonUnsupportedMediaType: "unsupported_media_type",
		ErrorReasonUnsupportedTarget:    "unsupported_target",
		ErrorReasonUnsupportedFeature:   "unsupported_feature",
		ErrorReasonDecodeFailed:         "decode_failed",
		ErrorReasonEncodeFailed:         "encode_failed",
		ErrorReasonInvalidDocument:      "invalid_document",
		ErrorReasonDepthExceeded:        "depth_exceeded",
		ErrorReasonInvalidNumber:        "invalid_number",
	}

	for reason, want := range testCases {
		if string(reason) != want {
			t.Fatalf("reason %q = %q; want %q", reason, reason, want)
		}
	}
}
