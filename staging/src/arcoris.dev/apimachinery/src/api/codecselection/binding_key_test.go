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

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestBindingKeysIncludeParameters(t *testing.T) {
	a := bindingRecord{
		contentType: testContentType(codec.MediaTypeJSON, MustParameter("profile", "public")),
		target:      codec.TargetObject,
		transport:   TransportBytes,
	}
	b := bindingRecord{
		contentType: testContentType(codec.MediaTypeJSON, MustParameter("profile", "storage")),
		target:      codec.TargetObject,
		transport:   TransportBytes,
	}

	if a.key() == b.key() {
		t.Fatalf("binding keys are equal; want parameters to distinguish keys")
	}
}
