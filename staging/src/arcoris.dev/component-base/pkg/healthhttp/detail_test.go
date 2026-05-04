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

package healthhttp

import "testing"

func TestDetailLevelString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  string
	}{
		{name: "none", level: DetailNone, want: "none"},
		{name: "failed", level: DetailFailed, want: "failed"},
		{name: "all", level: DetailAll, want: "all"},
		{name: "invalid", level: DetailLevel(99), want: "invalid"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.level.String(); got != test.want {
				t.Fatalf("String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestDetailLevelIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  bool
	}{
		{name: "none", level: DetailNone, want: true},
		{name: "failed", level: DetailFailed, want: true},
		{name: "all", level: DetailAll, want: true},
		{name: "invalid", level: DetailLevel(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.level.IsValid(); got != test.want {
				t.Fatalf("IsValid() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDetailLevelZeroValueIsNone(t *testing.T) {
	t.Parallel()

	var level DetailLevel
	if level != DetailNone {
		t.Fatalf("zero DetailLevel = %s, want %s", level, DetailNone)
	}
	if !level.IsValid() {
		t.Fatal("zero DetailLevel should be valid")
	}
	if level.IncludesChecks() {
		t.Fatal("zero DetailLevel should not include checks")
	}
}

func TestDetailLevelIncludesChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  bool
	}{
		{name: "none", level: DetailNone, want: false},
		{name: "failed", level: DetailFailed, want: true},
		{name: "all", level: DetailAll, want: true},
		{name: "invalid", level: DetailLevel(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.level.IncludesChecks(); got != test.want {
				t.Fatalf("IncludesChecks() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDetailLevelIncludesAllChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  bool
	}{
		{name: "none", level: DetailNone, want: false},
		{name: "failed", level: DetailFailed, want: false},
		{name: "all", level: DetailAll, want: true},
		{name: "invalid", level: DetailLevel(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.level.IncludesAllChecks(); got != test.want {
				t.Fatalf("IncludesAllChecks() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDetailLevelIncludesFailedChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  bool
	}{
		{name: "none", level: DetailNone, want: false},
		{name: "failed", level: DetailFailed, want: true},
		{name: "all", level: DetailAll, want: true},
		{name: "invalid", level: DetailLevel(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.level.IncludesFailedChecks(); got != test.want {
				t.Fatalf("IncludesFailedChecks() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestValidateDetailLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  bool
	}{
		{name: "none", level: DetailNone, want: true},
		{name: "failed", level: DetailFailed, want: true},
		{name: "all", level: DetailAll, want: true},
		{name: "invalid", level: DetailLevel(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validateDetailLevel(test.level)
			if got := err == nil; got != test.want {
				t.Fatalf("validateDetailLevel(%s) ok = %v, want %v; err=%v", test.level, got, test.want, err)
			}
		})
	}
}
