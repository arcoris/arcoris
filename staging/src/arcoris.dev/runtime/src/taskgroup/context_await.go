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

package taskgroup

import "context"

// AwaitContext blocks until ctx is done and returns the cancellation cause.
//
// If ctx has no explicit cause, AwaitContext falls back to ctx.Err.
// AwaitContext panics when ctx is nil. It is useful as a Task when a Group
// should include a context-driven sentinel task without defining additional
// runtime policy.
func AwaitContext(ctx context.Context) error {
	requireContext(ctx, errNilWaitContext)

	<-ctx.Done()
	return contextCause(ctx)
}
