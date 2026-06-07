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

package eval

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestInterruptedResult(t *testing.T) {
	t.Parallel()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 0)
	defer timeoutCancel()

	timeout := interruptedResult("storage", timeoutCtx)
	if timeout.Reason != health.ReasonTimeout || !errors.Is(timeout.Cause, context.DeadlineExceeded) {
		t.Fatalf("timeout result = %+v", timeout)
	}

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	canceled := interruptedResult("storage", canceledCtx)
	if canceled.Reason != health.ReasonCanceled || !errors.Is(canceled.Cause, context.Canceled) {
		t.Fatalf("canceled result = %+v", canceled)
	}

	active := interruptedResult("storage", context.Background())
	if active.Reason != health.ReasonCanceled || active.Cause != nil {
		t.Fatalf("active context result = %+v, want canceled reason with nil cause", active)
	}
}
