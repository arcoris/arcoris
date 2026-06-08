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

package memory

import "context"

// checkContext reports already-observed context cancellation.
//
// The memory store does not wait on context, start goroutines, or retry work.
// Context is accepted so the same Store contract can support future
// implementations that may cross process or storage boundaries.
func checkContext(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
