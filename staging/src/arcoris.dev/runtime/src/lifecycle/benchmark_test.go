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

package lifecycle

import (
	"context"
	"testing"
)

func BenchmarkLifecycleSnapshot(b *testing.B) {
	controller := NewController()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = controller.Snapshot()
	}
}

func BenchmarkLifecycleBeginStart(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		controller := NewController()
		_, _ = controller.BeginStart()
	}
}

func BenchmarkLifecycleFullStartStop(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		controller := NewController()
		_, _ = controller.BeginStart()
		_, _ = controller.MarkRunning()
		_, _ = controller.BeginStop()
		_, _ = controller.MarkStopped()
	}
}

func BenchmarkLifecycleWaitAlreadySatisfied(b *testing.B) {
	controller := NewController()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = controller.Wait(context.Background(), func(Snapshot) bool { return true })
	}
}
