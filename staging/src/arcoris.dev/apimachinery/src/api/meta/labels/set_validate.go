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

import "fmt"

// Validate checks every label key and value in the set.
func (s Set) Validate() error {
	for key, value := range s {
		if err := Key(key).Validate(); err != nil {
			return nested(fmt.Sprintf("labels[%q].key", key), ErrInvalidSet, err)
		}
		if err := Value(value).Validate(); err != nil {
			return nested(fmt.Sprintf("labels[%q].value", key), ErrInvalidSet, err)
		}
	}
	return nil
}
