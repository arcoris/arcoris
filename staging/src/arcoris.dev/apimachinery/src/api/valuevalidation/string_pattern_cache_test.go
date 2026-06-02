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

package valuevalidation

import "testing"

func TestCompilePatternCachesSuccessfulRegexp(t *testing.T) {
	run := newValidator(Options{})

	first, err := run.compilePattern(`^[a-z]+$`)
	if err != nil {
		t.Fatalf("compilePattern() error = %v", err)
	}
	second, err := run.compilePattern(`^[a-z]+$`)
	if err != nil {
		t.Fatalf("compilePattern() second error = %v", err)
	}

	if first != second {
		t.Fatalf("compilePattern() returned different regexp pointers")
	}
}

func TestCompilePatternCachesRegexpErrors(t *testing.T) {
	run := newValidator(Options{})

	_, firstErr := run.compilePattern(`[`)
	_, secondErr := run.compilePattern(`[`)

	if firstErr == nil {
		t.Fatalf("compilePattern() error = nil")
	}
	if firstErr != secondErr {
		t.Fatalf("compilePattern() returned different cached errors")
	}
}
