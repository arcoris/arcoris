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

package jsonconfig

import "testing"

func TestDateEncodingModeIsKnown(t *testing.T) {
	t.Parallel()

	if !isKnownDateEncodingMode(DateEncodingDefault) ||
		!isKnownDateEncodingMode(DateEncodingReject) ||
		!isKnownDateEncodingMode(DateEncodingISO8601) {
		t.Fatalf("known date encoding mode reported unknown")
	}
	if isKnownDateEncodingMode(DateEncodingMode(99)) {
		t.Fatalf("unknown date encoding mode reported known")
	}
}
