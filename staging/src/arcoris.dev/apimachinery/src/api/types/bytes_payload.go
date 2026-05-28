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

package types

// bytesPayload stores TypeBytes constraints.
type bytesPayload struct {
	// minLen is the inclusive minimum byte length.
	minLen intLimit
	// maxLen is the inclusive maximum byte length.
	maxLen intLimit
}

// cloneBytesPayload copies byte-sequence limits.
func cloneBytesPayload(p bytesPayload) bytesPayload {
	return p
}
