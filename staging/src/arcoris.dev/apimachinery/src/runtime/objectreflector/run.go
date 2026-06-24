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

package objectreflector

import (
	"context"
)

// Run reflects until ctx is canceled, a sink operation fails, or the source
// violates its contract.
//
// Run panics on a nil context. A nil context would create an uncancellable
// active runtime component, which is a programmer error.
func (r *Reflector) Run(ctx context.Context) error {
	if ctx == nil {
		panic("nil context")
	}
	if err := r.beginRun(); err != nil {
		return err
	}
	defer r.endRun()

	attempt := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := r.runCycle(ctx)
		if err == nil {
			attempt = 0
			continue
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if isRelistRequired(err) {
			attempt++
			if waitErr := r.waitBeforeRelist(ctx, attempt, err); waitErr != nil {
				return waitErr
			}
			continue
		}

		return err
	}
}

// waitBeforeRelist applies only to continuity-driven relists.
//
// Sink errors bypass RelistPolicy because the reflector has no repair or
// idempotency contract for partially applied sink writes.
func (r *Reflector) waitBeforeRelist(ctx context.Context, attempt int, cause error) error {
	return r.options.RelistPolicy.Wait(ctx, attempt, cause)
}
