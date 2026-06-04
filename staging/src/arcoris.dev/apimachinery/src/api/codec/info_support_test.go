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

func TestInfoSupportsTarget(t *testing.T) {
	info := Info{Targets: []Target{TargetObject, TargetValue}}

	if !info.Supports(TargetValue) {
		t.Fatalf("Supports(TargetValue) = false")
	}
	if info.Supports(TargetObjectOwnership) {
		t.Fatalf("Supports(TargetObjectOwnership) = true")
	}
}

func TestInfoSupportsMediaType(t *testing.T) {
	info := Info{MediaTypes: []MediaType{MediaTypeJSON, MediaTypeYAML}}

	if !info.SupportsMediaType(MediaTypeYAML) {
		t.Fatalf("SupportsMediaType(MediaTypeYAML) = false")
	}
	if info.SupportsMediaType(MediaTypeCBOR) {
		t.Fatalf("SupportsMediaType(MediaTypeCBOR) = true")
	}
}
