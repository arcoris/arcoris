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

package types

import "testing"

func FuzzParseFieldName(f *testing.F) {
	for _, seed := range []string{"name", "maxConcurrency", "", "Name", "1name", "naïve"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		name, err := ParseFieldName(s)
		if err == nil && !name.IsValid() {
			t.Fatalf("parsed invalid field name %q", s)
		}
	})
}

func FuzzParseTypeName(f *testing.F) {
	for _, seed := range []string{"meta.arcoris.dev.Name", "example.dev.CronExpression", "", "bad", "example.dev.name"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		name, err := ParseTypeName(s)
		if err == nil && !name.IsValid() {
			t.Fatalf("parsed invalid type name %q", s)
		}
	})
}

func FuzzValidateDescriptorDoesNotPanic(f *testing.F) {
	for _, seed := range []uint8{0, 1, 2, 23, 255} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, code uint8) {
		_ = ValidateLocal(Descriptor{code: DescriptorKind(code)})
	})
}

func FuzzValidateStringPatternDoesNotPanic(f *testing.F) {
	for _, seed := range []string{"^[a-z]+$", "[", "", ".*", `\pL+`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, pattern string) {
		_ = ValidateLocal(String().Pattern(pattern).Enum("alpha").Descriptor())
	})
}
