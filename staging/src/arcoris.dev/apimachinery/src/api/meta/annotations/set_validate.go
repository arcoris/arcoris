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

package annotations

import "fmt"

// ValidateLexical checks every annotation key and value in the set.
func (s Set) ValidateLexical() error {
	for key, value := range s {
		if err := key.ValidateLexical(); err != nil {
			return nested(fmt.Sprintf("annotations[%q].key", key.String()), ErrInvalidSet, err)
		}
		if err := value.ValidateLexical(); err != nil {
			return nested(fmt.Sprintf("annotations[%q].value", key.String()), ErrInvalidSet, err)
		}
	}
	return nil
}
