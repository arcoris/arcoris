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

func BenchmarkRecordViewGet(b *testing.B) {
	value := benchmarkNestedRecordValue()
	view, ok := value.AsRecord()
	if !ok {
		b.Fatal("expected record view")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink, benchmarkOKSink = view.Get("payload")
	}
}

func BenchmarkListViewAt(b *testing.B) {
	value := MustListValue(
		StringValue("first"),
		benchmarkNestedRecordValue(),
		BytesValue([]byte("payload")),
	)
	view, ok := value.AsList()
	if !ok {
		b.Fatal("expected list view")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink, benchmarkOKSink = view.At(1)
	}
}

func BenchmarkValueCloneNestedRecord(b *testing.B) {
	value := benchmarkNestedRecordValue()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkValueSink = value.Clone()
	}
}

// benchmarkNestedRecordValue returns a small but nested payload for clone and
// view benchmarks.
func benchmarkNestedRecordValue() Value {
	return MustRecordValue(
		MustRecordMember("name", StringValue("worker")),
		MustRecordMember("payload", BytesValue([]byte("payload"))),
		MustRecordMember("tags", MustListValue(
			StringValue("control"),
			StringValue("active"),
			MustRecordValue(
				MustRecordMember("zone", StringValue("east")),
				MustRecordMember("tier", StringValue("primary")),
			),
		)),
	)
}
