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

package capacity

import "arcoris.dev/capacity/internal/diagnostic"

// errorAt builds one direct capacity diagnostic.
func errorAt(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// panicAt panics with one structured capacity diagnostic.
func panicAt(path string, err error, reason ErrorReason, detail string) {
	panic(errorAt(path, err, reason, detail))
}
