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

import (
	"errors"
	"testing"
)

func TestValidatePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "startup", path: DefaultStartupPath, want: true},
		{name: "live", path: DefaultLivePath, want: true},
		{name: "ready", path: DefaultReadyPath, want: true},
		{name: "healthz", path: DefaultHealthPath, want: true},
		{name: "health", path: DefaultHealthPlainPath, want: true},
		{name: "nested", path: "/internal/health/ready", want: true},
		{name: "empty", path: "", want: false},
		{name: "relative", path: "readyz", want: false},
		{name: "root", path: "/", want: false},
		{name: "query", path: "/readyz?verbose", want: false},
		{name: "fragment", path: "/readyz#fragment", want: false},
		{name: "absolute_url", path: "http://localhost/readyz", want: false},
		{name: "protocol_relative", path: "//localhost/readyz", want: false},
		{name: "scheme_inside", path: "/http://localhost/readyz", want: false},
		{name: "space", path: "/readyz debug", want: false},
		{name: "tab", path: "/readyz\tdebug", want: false},
		{name: "newline", path: "/readyz\ndebug", want: false},
		{name: "backslash", path: "/readyz\\debug", want: false},
		{name: "delete_control", path: "/readyz" + string(rune(0x7f)), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidatePath(test.path)
			if got := err == nil; got != test.want {
				t.Fatalf("ValidatePath(%q) ok = %v, want %v; err=%v", test.path, got, test.want, err)
			}
			if test.want {
				return
			}
			if !errors.Is(err, ErrInvalidPath) {
				t.Fatalf("ValidatePath(%q) error = %v, want ErrInvalidPath", test.path, err)
			}

			var pathErr InvalidPathError
			if !errors.As(err, &pathErr) {
				t.Fatalf("ValidatePath(%q) error = %T, want InvalidPathError", test.path, err)
			}
			if pathErr.Path != test.path {
				t.Fatalf("InvalidPathError.Path = %q, want %q", pathErr.Path, test.path)
			}
		})
	}
}

func TestValidPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "valid", path: "/readyz", want: true},
		{name: "invalid", path: "readyz", want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := validPath(test.path); got != test.want {
				t.Fatalf("validPath(%q) = %v, want %v", test.path, got, test.want)
			}
		})
	}
}

func TestInvalidPathRune(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{name: "letter", r: 'a', want: false},
		{name: "slash", r: '/', want: false},
		{name: "space", r: ' ', want: true},
		{name: "tab", r: '\t', want: true},
		{name: "newline", r: '\n', want: true},
		{name: "nul", r: 0x00, want: true},
		{name: "delete", r: 0x7f, want: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := invalidPathRune(test.r); got != test.want {
				t.Fatalf("invalidPathRune(%q) = %v, want %v", test.r, got, test.want)
			}
		})
	}
}
