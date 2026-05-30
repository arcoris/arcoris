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

package stamp

import "strconv"

// Generation tracks desired-state changes for one object.
//
// Zero means absent or unassigned. Generation is not a resource version and not
// a storage revision.
type Generation uint64

// IsZero reports whether the generation is absent.
func (g Generation) IsZero() bool {
	return g == 0
}

// String returns the decimal generation text.
func (g Generation) String() string {
	return strconv.FormatUint(uint64(g), 10)
}
