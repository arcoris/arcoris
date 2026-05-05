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

func TestDefaultPathConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "startup", got: DefaultStartupPath, want: "/startupz"},
		{name: "live", got: DefaultLivePath, want: "/livez"},
		{name: "ready", got: DefaultReadyPath, want: "/readyz"},
		{name: "healthz", got: DefaultHealthPath, want: "/healthz"},
		{name: "health", got: DefaultHealthPlainPath, want: "/health"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.got != test.want {
				t.Fatalf("path constant = %q, want %q", test.got, test.want)
			}
		})
	}
}

func TestPrimaryDefaultPathsAreValid(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		DefaultStartupPath,
		DefaultLivePath,
		DefaultReadyPath,
	} {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			if err := ValidatePath(path); err != nil {
				t.Fatalf("ValidatePath(%q) = %v, want nil", path, err)
			}
		})
	}
}

func TestCompatibilityPathsAreValid(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		DefaultHealthPath,
		DefaultHealthPlainPath,
	} {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			if err := ValidatePath(path); err != nil {
				t.Fatalf("ValidatePath(%q) = %v, want nil", path, err)
			}
		})
	}
}
