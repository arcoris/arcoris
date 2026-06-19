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
	"testing"
)

func TestRunRejectsConcurrentCalls(t *testing.T) {
	stream := waitingStream()
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: testRead(t, 0)}},
		watchResponses: []watchResponse{{stream: stream}},
	}
	reflector := newTestReflector(t, source, newRecordingSink(1))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)

	go func() {
		done <- reflector.Run(ctx)
	}()
	<-stream.nextStarted

	err := reflector.Run(context.Background())
	requireErrorIs(t, err, ErrAlreadyRunning)

	cancel()
	requireErrorIs(t, <-done, context.Canceled)
}

func TestRunMayStartAgainAfterExit(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	requireErrorIs(t, reflector.Run(ctx), context.Canceled)
	requireErrorIs(t, reflector.Run(ctx), context.Canceled)
}
