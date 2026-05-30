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

package labels

import "arcoris.dev/apimachinery/api/meta/internal/metagrammar"

// maxLabelValueLength keeps labels compact enough for indexing-oriented metadata.
const maxLabelValueLength = 63

// Validate checks the label value grammar.
func (v Value) Validate() error {
	return fromGrammar(
		"label.value",
		ErrInvalidValue,
		metagrammar.ValidateMapValue(v.String(), metagrammar.MapValueOptions{
			AllowEmpty: true,
			MaxLength:  maxLabelValueLength,
			Strict:     true,
		}),
	)
}
