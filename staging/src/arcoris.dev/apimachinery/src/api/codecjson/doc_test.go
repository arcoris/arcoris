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

package codecjson

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

// TestPackageDocumentedTargets checks the targets promised by package docs.
func TestPackageDocumentedTargets(t *testing.T) {
	info := newTestCodec(t).Info()

	for _, target := range []codec.Target{
		codec.TargetValue,
		codec.TargetObject,
		codec.TargetObjectOwnership,
	} {
		if !info.Supports(target) {
			t.Fatalf("Info().Supports(%q) = false", target)
		}
	}
}
