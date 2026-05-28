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

// float64Limit stores optional float64 constraints without pointers.
//
// Floating-point min/max rules need the same unset-vs-zero distinction as int
// rules. The wrapper also gives validation one place to reject NaN and
// infinities without making builder methods allocate.
type float64Limit struct {
	// value is meaningful only when set is true.
	//
	// Validation code must check set before treating value as an active rule.
	value float64
	// set reports whether value was explicitly configured.
	//
	// value=0 with set=false means no configured rule; value=0 with set=true
	// means the rule is explicitly zero.
	set bool
}
