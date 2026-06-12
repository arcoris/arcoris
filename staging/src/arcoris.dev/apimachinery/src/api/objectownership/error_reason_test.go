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

package objectownership

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	tests := []struct {
		name string
		got  ErrorReason
		want string
	}{
		{name: "invalid state", got: ErrorReasonInvalidState, want: "invalid_state"},
		{name: "invalid desired", got: ErrorReasonInvalidDesired, want: "invalid_desired"},
		{name: "invalid observed", got: ErrorReasonInvalidObserved, want: "invalid_observed"},
		{name: "invalid metadata labels", got: ErrorReasonInvalidMetadataLabels, want: "invalid_metadata_labels"},
		{name: "invalid metadata annotations", got: ErrorReasonInvalidMetadataAnnotations, want: "invalid_metadata_annotations"},
		{name: "not normalized", got: ErrorReasonNotNormalized, want: "not_normalized"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Fatalf("reason = %q; want %q", tt.got, tt.want)
			}
		})
	}
}
