/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package bulkhead

import (
	"arcoris.dev/admission"
	"arcoris.dev/snapshot"
)

var (
	// Compile-time contract checks for Bulkhead's state snapshot APIs.
	_ snapshot.Source[Snapshot] = (*Bulkhead)(nil)
	_ snapshot.RevisionSource   = (*Bulkhead)(nil)

	// Compile-time contract check for the admission-compatible non-blocking API.
	//
	// This does not make Bulkhead a general admission controller. It only states
	// that the local bounded in-flight primitive can expose its existing
	// check-and-reserve operation through admission's generic Result contract.
	_ admission.Admitter[Request, *Lease, snapshot.Snapshot[Snapshot]] = (*Bulkhead)(nil)
)
