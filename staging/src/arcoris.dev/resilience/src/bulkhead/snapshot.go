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

package bulkhead

import "arcoris.dev/capacity"

// Snapshot is the copyable state model of one Bulkhead's local capacity.
//
// Snapshot is an alias of capacity.ScalarSnapshot because bulkhead does not maintain
// parallel diagnostic counters. Limit, Reserved, Available, and Debt are exactly
// the live in-flight capacity state owned by the underlying ledger. This keeps
// debt semantics identical when a limit is reduced below active leases.
type Snapshot = capacity.ScalarSnapshot
