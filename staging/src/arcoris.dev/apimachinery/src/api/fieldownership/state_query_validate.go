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

// validateOwnedPathQuery checks one caller-supplied ownership query path.
func validateOwnedPathQuery(path fieldpath.Path, detail string) error {
	if err := path.ValidateStructure(); err != nil {
		return wrapAt(
			"owned."+path.CanonicalText(),
			ErrInvalidPath,
			ErrorReasonInvalidOwnedPath,
			detail,
			err,
		)
	}

	return nil
}
