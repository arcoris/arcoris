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

package identity

import "testing"

var (
	benchmarkGVRPSink   GroupVersionResourcePath
	benchmarkStringSink string
	benchmarkBytesSink  []byte
)

func BenchmarkParseGroupVersionResourcePath(b *testing.B) {
	const input = "control.arcoris.dev/v1:workers/status"

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parsed, err := ParseGroupVersionResourcePath(input)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkGVRPSink = parsed
	}
}

func BenchmarkGroupVersionResourcePathValidate(b *testing.B) {
	identity := GroupVersionResourcePath{
		Group:       Group("control.arcoris.dev"),
		Version:     Version("v1"),
		Resource:    Resource("workers"),
		Subresource: Subresource("status"),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := identity.Validate(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGroupVersionResourcePathString(b *testing.B) {
	identity := GroupVersionResourcePath{
		Group:       Group("control.arcoris.dev"),
		Version:     Version("v1"),
		Resource:    Resource("workers"),
		Subresource: Subresource("status"),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkStringSink = identity.String()
	}
}

func BenchmarkGroupVersionResourcePathMarshalJSON(b *testing.B) {
	identity := GroupVersionResourcePath{
		Group:       Group("control.arcoris.dev"),
		Version:     Version("v1"),
		Resource:    Resource("workers"),
		Subresource: Subresource("status"),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := identity.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytesSink = data
	}
}
