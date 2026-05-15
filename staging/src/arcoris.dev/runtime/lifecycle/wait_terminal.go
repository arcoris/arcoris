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

package lifecycle

import "context"

// WaitTerminal blocks until the lifecycle reaches a terminal state or ctx is
// cancelled.
//
// Terminal states are StateStopped and StateFailed. The returned Snapshot lets
// callers distinguish successful shutdown from failure:
//
//	snapshot, err := controller.WaitTerminal(ctx)
//	if err != nil {
//		return err
//	}
//	if snapshot.IsFailed() {
//		return snapshot.FailureCause
//	}
//
// WaitTerminal is safe to call concurrently with transition methods.
func (c *Controller) WaitTerminal(ctx context.Context) (Snapshot, error) {
	return c.Wait(ctx, func(snap Snapshot) bool {
		return snap.IsTerminal()
	})
}
