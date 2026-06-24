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
	"time"
)

// RelistPolicy paces the next list-watch cycle after continuity is explicitly
// lost.
//
// attempt counts consecutive relist-required cycle endings within one Run call.
// The first relist-required cycle uses attempt=1. cause is the error that ended
// the previous cycle. Returning nil lets Run start the next ListCollection
// cycle; returning a non-nil error makes Run return that error.
type RelistPolicy interface {
	Wait(ctx context.Context, attempt int, cause error) error
}

// ImmediateRelistPolicy preserves the historical reflector behavior: relist as
// soon as the current context still permits work.
type ImmediateRelistPolicy struct{}

// Wait returns ctx.Err and otherwise performs no delay.
func (ImmediateRelistPolicy) Wait(ctx context.Context, _ int, _ error) error {
	return ctx.Err()
}

// ConstantRelistDelay waits a fixed duration before the next list-watch cycle.
//
// This is intentionally small and standard-library-only. It is not exponential
// backoff, jitter, or a generic retry policy.
type ConstantRelistDelay struct {
	Delay time.Duration
}

// Wait blocks for Delay or context cancellation. Delay <= 0 is immediate.
func (p ConstantRelistDelay) Wait(ctx context.Context, _ int, _ error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if p.Delay <= 0 {
		return nil
	}

	timer := time.NewTimer(p.Delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
