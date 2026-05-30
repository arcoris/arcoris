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

package owner

import "fmt"

// Validate checks owner references and enforces owner-list invariants.
func (l List) Validate() error {
	controllers := 0
	seen := make(map[Reference]struct{}, len(l))

	for i, ref := range l {
		if err := ref.Validate(); err != nil {
			return nested(fmt.Sprintf("ownerReferences[%d]", i), ErrInvalidList, err)
		}
		if ref.Controller {
			controllers++
			if controllers > 1 {
				return invalid(
					fmt.Sprintf("ownerReferences[%d].controller", i),
					ErrMultipleControllers,
					ErrorReasonMultipleControllers,
					"at most one owner reference may be marked as controller",
				)
			}
		}
		if _, ok := seen[ref]; ok {
			return invalid(
				fmt.Sprintf("ownerReferences[%d]", i),
				ErrDuplicateReference,
				ErrorReasonDuplicateReference,
				"owner references must be unique",
			)
		}
		seen[ref] = struct{}{}
	}
	return nil
}
