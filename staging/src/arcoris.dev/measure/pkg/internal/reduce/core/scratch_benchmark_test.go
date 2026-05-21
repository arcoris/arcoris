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

package core

import "testing"

func BenchmarkScratchReset(b *testing.B) {
	s := scratchForBenchmark()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Reset()
	}
}

func BenchmarkScratchClear(b *testing.B) {
	s := scratchForBenchmark()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Clear()
	}
}

func BenchmarkScratchRelease(b *testing.B) {
	s := scratchForBenchmark()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Release()
		s = scratchForBenchmark()
	}
}

func BenchmarkEnsurePartials(b *testing.B) {
	var s Scratch[int]
	s.EnsurePartialsDirty(1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = s.EnsurePartials(1024)
	}
}

func BenchmarkEnsurePartialsDirty(b *testing.B) {
	var s Scratch[int]
	s.EnsurePartialsDirty(1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = s.EnsurePartialsDirty(1024)
	}
}

func BenchmarkEnsureUsed(b *testing.B) {
	var s Scratch[int]
	s.EnsureUsedDirty(1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = s.EnsureUsed(1024)
	}
}

func BenchmarkEnsureUsedDirty(b *testing.B) {
	var s Scratch[int]
	s.EnsureUsedDirty(1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = s.EnsureUsedDirty(1024)
	}
}

func scratchForBenchmark() Scratch[*int] {
	return Scratch[*int]{
		Ranges:   make([]Range, 1024),
		Partials: make([]*int, 1024),
		Used:     make([]bool, 1024),
	}
}
