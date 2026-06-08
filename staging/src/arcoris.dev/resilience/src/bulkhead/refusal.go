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

// Refusal classifies a direct bulkhead acquisition refusal.
//
// Refusal is an alias of capacity.Refusal. Bulkhead does not duplicate scalar
// accounting taxonomy; it gives capacity refusals resilience-domain meaning for
// protected in-flight work acquisition.
//
// The common scalar outcomes are RefusalNone, RefusalInsufficient, and
// RefusalDebt. Keeping the alias lets callers use capacity's validation and
// string semantics without introducing a second almost-identical taxonomy in the
// resilience layer.
type Refusal = capacity.Refusal

const (
	// RefusalNone means local accounting did not refuse the attempt.
	RefusalNone Refusal = capacity.RefusalNone

	// RefusalInsufficient means the bulkhead had too little available capacity.
	RefusalInsufficient Refusal = capacity.RefusalInsufficient

	// RefusalDebt means active leases already exceed the current limit.
	RefusalDebt Refusal = capacity.RefusalDebt

	// RefusalUnknownResource is preserved for callers that handle the complete
	// capacity refusal taxonomy. Scalar bulkheads do not normally produce it.
	RefusalUnknownResource Refusal = capacity.RefusalUnknownResource
)
