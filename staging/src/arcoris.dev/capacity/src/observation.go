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

import "arcoris.dev/snapshot"

// Observation is the diagnostic read model returned by scalar observed methods.
type Observation struct {
	// Snapshot is read after the observed attempt.
	Snapshot snapshot.Snapshot[Snapshot]

	// Refusal classifies why an observed scalar attempt did not reserve capacity.
	Refusal Refusal
}
