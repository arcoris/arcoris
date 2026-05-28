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

// uint64Limit stores optional uint64 constraints without pointers.
//
// Descriptor rules need to distinguish "unset" from "set to zero". A pointer
// would also encode that distinction, but it would add hidden allocations to
// common builder paths. The explicit set bit keeps the descriptor compact,
// allocation-light, and obvious during validation.
type uint64Limit struct {
	// value is meaningful only when set is true.
	//
	// Callers must never read value as an active constraint without checking set.
	value uint64
	// set reports whether value was explicitly configured.
	//
	// A zero value with set=false means "no uint64 limit"; value=0 with set=true
	// means "the configured limit is exactly zero".
	set bool
}
