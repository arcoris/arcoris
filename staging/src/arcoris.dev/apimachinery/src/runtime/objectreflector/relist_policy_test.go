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
	"errors"
	"sync"
	"testing"
	"time"
)

type relistCall struct {
	attempt int
	cause   error
}

type recordingRelistPolicy struct {
	mu sync.Mutex

	calls []relistCall
	err   error
}

func (p *recordingRelistPolicy) Wait(_ context.Context, attempt int, cause error) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.calls = append(p.calls, relistCall{attempt: attempt, cause: cause})

	return p.err
}

func (p *recordingRelistPolicy) recordedCalls() []relistCall {
	p.mu.Lock()
	defer p.mu.Unlock()

	return append([]relistCall(nil), p.calls...)
}

type blockingRelistPolicy struct {
	entered chan struct{}
	once    sync.Once
}

func newBlockingRelistPolicy() *blockingRelistPolicy {
	return &blockingRelistPolicy{entered: make(chan struct{})}
}

func (p *blockingRelistPolicy) Wait(ctx context.Context, _ int, _ error) error {
	p.once.Do(func() { close(p.entered) })
	<-ctx.Done()
	return ctx.Err()
}

func TestImmediateRelistPolicyReturnsContextErrorOnly(t *testing.T) {
	requireNoError(t, ImmediateRelistPolicy{}.Wait(context.Background(), 1, errors.New("gap")))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	requireErrorIs(t, ImmediateRelistPolicy{}.Wait(ctx, 1, errors.New("gap")), context.Canceled)
}

func TestConstantRelistDelayImmediateForNonPositiveDelay(t *testing.T) {
	requireNoError(t, ConstantRelistDelay{}.Wait(context.Background(), 1, errors.New("gap")))
	requireNoError(t, ConstantRelistDelay{Delay: -time.Second}.Wait(context.Background(), 1, errors.New("gap")))
}

func TestConstantRelistDelayReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := (ConstantRelistDelay{Delay: time.Hour}).Wait(ctx, 1, errors.New("gap"))

	requireErrorIs(t, err, context.Canceled)
}
