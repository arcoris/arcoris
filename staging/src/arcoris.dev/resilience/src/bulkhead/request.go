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

// Request is the typed admission request consumed by Bulkhead.TryAdmit.
//
// TryAcquire is the ergonomic direct API for reserving one in-flight unit.
// TryAdmit uses Request so Bulkhead can satisfy
// admission.Admitter[Request, *Lease, snapshot.Snapshot[Snapshot]] without
// turning this package into a generic admission framework.
//
// Amount is the number of local in-flight capacity units to reserve. It uses the
// same scalar type as capacity.Ledger through the bulkhead.Amount alias. Request
// intentionally contains no context, priority, tenant, request class, metadata,
// deadline, or queueing policy; those dimensions belong to higher-level
// admission or policy layers.
//
// Invalid Amount values remain programming or configuration errors handled by
// the underlying capacity validation path. They are not converted into denied
// admission results because denial represents ordinary back-pressure, not a
// malformed request.
type Request struct {
	// Amount is the number of in-flight capacity units requested by TryAdmit.
	Amount Amount
}
