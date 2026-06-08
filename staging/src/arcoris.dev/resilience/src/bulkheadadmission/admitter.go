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

package bulkheadadmission

import "arcoris.dev/resilience/bulkhead"

// Admitter adapts a Bulkhead to admission.Result.
//
// Admitter does not add synchronization or accounting of its own. The wrapped
// Bulkhead remains the sole owner of scalar capacity and lease release
// lifecycle. Admitter is safe to copy by value because it only carries the
// pointer to the underlying Bulkhead; lease ownership still belongs to the grant
// returned from a successful TryAdmit call.
type Admitter struct {
	// bulkhead is the direct non-blocking primitive being adapted.
	//
	// It may be nil only as programmer misuse. TryAdmit intentionally lets the
	// core bulkhead receiver validation produce the package-owned panic so the
	// direct and adapter surfaces share one misuse taxonomy.
	bulkhead *bulkhead.Bulkhead
}

// New returns an admission adapter for b.
//
// Nil and uninitialized bulkheads are rejected by the core bulkhead methods when
// TryAdmit is called. New performs no eager validation so adapters can be stored
// in fixtures and initialized by the same path as direct callers. The adapter
// does not take ownership of b; callers must keep the Bulkhead alive for the
// adapter's lifetime.
func New(b *bulkhead.Bulkhead) Admitter {
	return Admitter{bulkhead: b}
}
