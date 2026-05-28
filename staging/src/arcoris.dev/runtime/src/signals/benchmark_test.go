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

package signals

import (
	"context"
	"os"
	"testing"
)

func BenchmarkSignalSetUnique(b *testing.B) {
	sigs := []os.Signal{testSIGINT, testSIGTERM, testSIGINT, testSIGHUP}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Unique(sigs)
	}
}

func BenchmarkSignalSetContains(b *testing.B) {
	sigs := []os.Signal{testSIGINT, testSIGTERM, testSIGHUP}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Contains(sigs, testSIGTERM)
	}
}

func BenchmarkSignalCause(b *testing.B) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(NewSignalError(testSIGTERM))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = Cause(ctx)
	}
}
