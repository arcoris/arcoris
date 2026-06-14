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

import "arcoris.dev/apimachinery/api/meta/labels"

// validate checks label requirement structure and delegates label lexical rules
// to api/meta/labels.
func (r LabelRequirement) validate(path string) error {
	return r.req.validate(path, validateLabelKey, validateLabelValue)
}

// validateLabelKey delegates label key syntax to api/meta/labels.
func validateLabelKey(key string) error {
	return labels.Key(key).ValidateLexical()
}

// validateLabelValue delegates label value syntax to api/meta/labels.
func validateLabelValue(value string) error {
	return labels.Value(value).ValidateLexical()
}
