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
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/value"
)

// preserveWithoutOverlayContainer preserves base when neither side can expose
// the selected descendants under expected.
func preserveWithoutOverlayContainer(
	base operand,
	overlay operand,
	expected value.Kind,
) (operand, bool) {
	if hasKind(base, expected) || hasKind(overlay, expected) {
		return operand{}, false
	}
	if base.Present() {
		return base.Clone(), true
	}

	return valuepresence.Absent(), true
}

// hasKind reports whether o is present and stores kind.
func hasKind(o operand, kind value.Kind) bool {
	return o.Present() && o.Value().Kind() == kind
}
