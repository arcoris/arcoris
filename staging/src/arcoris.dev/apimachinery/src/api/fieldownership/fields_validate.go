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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// validateFields defensively validates every path contained in fields.
//
// Public fieldpath.Set constructors already validate paths. This extra check
// keeps State robust if a future internal caller receives a malformed set value.
func validateFields(fields fieldpath.Set, detail string) error {
	return validateFieldsAt("fields", fields, ErrorReasonInvalidPath, detail)
}

// validateFieldsAt defensively validates every path contained in fields.
func validateFieldsAt(path string, fields fieldpath.Set, reason ErrorReason, detail string) error {
	for _, p := range fields.Paths() {
		if err := p.ValidateStructure(); err != nil {
			return wrapAt(
				path+"."+p.CanonicalText(),
				ErrInvalidPath,
				reason,
				detail,
				err,
			)
		}
	}

	return nil
}
