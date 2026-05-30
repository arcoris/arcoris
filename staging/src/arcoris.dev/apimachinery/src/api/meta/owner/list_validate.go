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

import (
	"fmt"

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

// Validate checks owner references and enforces owner-list invariants.
func (l List) Validate() error {
	controllerCount := 0
	seenRefs := make(map[metaidentity.ObjectReference]struct{}, len(l))

	for index, reference := range l {
		path := fmt.Sprintf("ownerReferences[%d]", index)

		if err := reference.Validate(); err != nil {
			return nested(path, ErrInvalidList, err)
		}

		if reference.Controller {
			controllerCount++
			if controllerCount > 1 {
				return invalid(
					path+".controller",
					ErrMultipleControllers,
					ErrorReasonMultipleControllers,
					"at most one owner reference may be marked as controller",
				)
			}
		}

		if _, ok := seenRefs[reference.Ref]; ok {
			return invalid(
				path,
				ErrDuplicateReference,
				ErrorReasonDuplicateReference,
				"owner references must be unique",
			)
		}

		seenRefs[reference.Ref] = struct{}{}
	}

	return nil
}
