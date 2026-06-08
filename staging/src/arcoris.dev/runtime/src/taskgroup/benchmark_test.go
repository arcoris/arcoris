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

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func BenchmarkGroupWaitNoTasks(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		group := New(context.Background())
		_ = group.Wait()
	}
}

func BenchmarkGroupGoWaitNoopTask(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		group := New(context.Background())
		group.Go("task", func(context.Context) error { return nil })
		_ = group.Wait()
	}
}

func BenchmarkGroupGoWaitManyNoopTasks(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		group := New(context.Background())
		for i := 0; i < 8; i++ {
			group.Go(fmt.Sprintf("task-%d", i), func(context.Context) error { return nil })
		}
		_ = group.Wait()
	}
}

func BenchmarkGroupCancel(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		group := New(context.Background())
		group.Cancel(nil)
	}
}

func BenchmarkGroupTaskErrorsJoin(b *testing.B) {
	err := errors.Join(
		TaskError{Name: "a", Err: errors.New("a")},
		TaskError{Name: "b", Err: errors.New("b")},
		TaskError{Name: "c", Err: errors.New("c")},
	)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = TaskErrors(err)
	}
}
