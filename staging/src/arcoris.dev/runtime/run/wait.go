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

package run

import "context"

// Wait blocks until ctx is done and returns the context cancellation cause.
//
// If ctx has no explicit cause, Wait falls back to ctx.Err. Wait panics when ctx
// is nil. Wait is useful as a Task when a Group should include a context-driven
// sentinel task.
func Wait(ctx context.Context) error {
	requireContext(ctx, errNilWaitContext)

	<-ctx.Done()
	return contextCause(ctx)
}
