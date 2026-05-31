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

package value

import "testing"

var (
	benchmarkValueSink Value
	benchmarkOKSink    bool
)

func BenchmarkObjectViewGet(b *testing.B) {
	value := benchmarkNestedObjectValue()
	view, ok := value.Object()
	if !ok {
		b.Fatal("expected object view")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink, benchmarkOKSink = view.Get("payload")
	}
}

func BenchmarkListViewAt(b *testing.B) {
	value := MustList(
		String("first"),
		benchmarkNestedObjectValue(),
		Bytes([]byte("payload")),
	)
	view, ok := value.List()
	if !ok {
		b.Fatal("expected list view")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink, benchmarkOKSink = view.At(1)
	}
}

func BenchmarkMapViewGet(b *testing.B) {
	value := MustMap(
		MapEntry("name", String("worker")),
		MapEntry("payload", benchmarkNestedObjectValue()),
		MapEntry("token", Bytes([]byte("token"))),
	)
	view, ok := value.Map()
	if !ok {
		b.Fatal("expected map view")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink, benchmarkOKSink = view.Get("payload")
	}
}

func BenchmarkValueCloneNestedObject(b *testing.B) {
	value := benchmarkNestedObjectValue()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink = value.Clone()
	}
}

// benchmarkNestedObjectValue returns a small but nested payload for clone and
// view benchmarks.
func benchmarkNestedObjectValue() Value {
	return MustObject(
		ObjectField("name", String("worker")),
		ObjectField("payload", Bytes([]byte("payload"))),
		ObjectField("tags", MustList(
			String("control"),
			String("active"),
			MustMap(
				MapEntry("zone", String("east")),
				MapEntry("tier", String("primary")),
			),
		)),
	)
}
