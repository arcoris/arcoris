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

package diagnostic

import "errors"

// Unwrap preserves the broad sentinel and nested cause identities.
func (r Record[R]) Unwrap() error {
	switch {
	case r.Err != nil && r.Cause != nil:
		return errors.Join(r.Err, r.Cause)
	case r.Err != nil:
		return r.Err
	default:
		return r.Cause
	}
}
