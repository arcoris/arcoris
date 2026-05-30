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

package meta

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
)

func TestTypeMetaClone(t *testing.T) {
	meta := TypeMeta{
		APIVersion: apiidentity.GroupVersion{Group: "control.arcoris.dev", Version: "v1"},
		Kind:       "Worker",
	}

	if meta.Clone() != meta {
		t.Fatal("Clone() changed value")
	}
}
