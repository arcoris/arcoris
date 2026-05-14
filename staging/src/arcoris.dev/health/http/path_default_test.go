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
	"testing"

	"arcoris.dev/health"
)

func TestDefaultPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target health.Target
		path   string
		ok     bool
	}{
		{name: "startup", target: health.TargetStartup, path: DefaultStartupPath, ok: true},
		{name: "live", target: health.TargetLive, path: DefaultLivePath, ok: true},
		{name: "ready", target: health.TargetReady, path: DefaultReadyPath, ok: true},
		{name: "unknown", target: health.TargetUnknown, ok: false},
		{name: "invalid", target: health.Target(99), ok: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			path, ok := DefaultPath(test.target)
			if ok != test.ok {
				t.Fatalf("DefaultPath(%s) ok = %v, want %v", test.target, ok, test.ok)
			}
			if path != test.path {
				t.Fatalf("DefaultPath(%s) path = %q, want %q", test.target, path, test.path)
			}
		})
	}
}

func TestCompatibilityPathsAreNotTargetDefaults(t *testing.T) {
	t.Parallel()

	for _, target := range []health.Target{
		health.TargetStartup,
		health.TargetLive,
		health.TargetReady,
	} {
		target := target
		t.Run(target.String(), func(t *testing.T) {
			t.Parallel()

			path, ok := DefaultPath(target)
			if !ok {
				t.Fatalf("DefaultPath(%s) ok = false, want true", target)
			}
			if path == DefaultHealthPath || path == DefaultHealthPlainPath {
				t.Fatalf(
					"DefaultPath(%s) = %q, compatibility paths must not be target defaults",
					target,
					path,
				)
			}
		})
	}
}
