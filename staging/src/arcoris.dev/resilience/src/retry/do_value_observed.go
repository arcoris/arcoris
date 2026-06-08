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

package retry

import "context"

// DoValueObserved executes op with bounded retry orchestration and returns the
// successful value plus terminal Outcome.
//
// DoValueObserved is the observed form of DoValue. It exposes completion
// metadata directly without requiring an Observer. The returned Outcome is valid
// for every non-panic terminal path and matches the Outcome delivered in the
// terminal EventRetryStop observer event.
//
// The returned value is meaningful only when err is nil. Failed-attempt values
// are ignored, and terminal failures return the zero value of T with the
// terminal Outcome and error.
//
// DoValueObserved does not recover panics from op, classifiers, observers,
// clocks, or delay schedules. Panics mean no retry-owned terminal Outcome is
// produced.
func DoValueObserved[T any](
	ctx context.Context,
	op ValueOperation[T],
	opts ...Option,
) (T, Outcome, error) {
	requireContext(ctx)
	requireValueOperation(op)

	return run(ctx, op, configOf(opts...))
}
