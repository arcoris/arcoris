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

func TestMediaTypeString(t *testing.T) {
	if MediaTypeJSON.String() != "application/json" {
		t.Fatalf("MediaType.String() = %q", MediaTypeJSON.String())
	}
}

func TestMediaTypeIsZero(t *testing.T) {
	if !MediaType("").IsZero() {
		t.Fatalf("zero media type IsZero() = false")
	}
	if MediaTypeJSON.IsZero() {
		t.Fatalf("json media type IsZero() = true")
	}
}
