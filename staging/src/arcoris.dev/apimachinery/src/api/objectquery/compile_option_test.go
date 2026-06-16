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

package objectquery

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestNewCompileOptionsIgnoresNilOptionsAndAppliesFields verifies option
// assembly stays convenient for callers that build option slices conditionally.
func TestNewCompileOptionsIgnoresNilOptionsAndAppliesFields(t *testing.T) {
	fields := mustFieldSet(t, selectable(fieldRef("spec.image"), value.KindString, Operators(OperatorEquals)))

	opts := newCompileOptions([]CompileOption{nil, WithSelectableFields(fields)})

	if opts.fields == nil {
		t.Fatal("opts.fields = nil; want selectable field set")
	}
}
