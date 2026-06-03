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

package valuemerge

import (
	"testing"

	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestRequireValidValueRejectsZero(t *testing.T) {
	err := requireValidValue(root(), valuepresence.Present(value.Value{}))

	requireErrorIs(t, err, ErrInvalidValue)
}

func TestRequireKindRejectsMismatch(t *testing.T) {
	err := requireKind(root(), valuepresence.Present(str("x")), value.KindBool)

	requireErrorIs(t, err, ErrKindMismatch)
}

func TestScalarKindIncludesNull(t *testing.T) {
	got, ok := scalarKind(types.TypeNull)

	if !ok {
		t.Fatalf("scalarKind(TypeNull) ok = false")
	}
	if got != value.KindNull {
		t.Fatalf("kind = %s; want null", got)
	}
}
