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

package healthgrpc

import (
	"errors"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/health/eval"
	"arcoris.dev/health/healthtest"
)

func TestNewServerRejectsNilSource(t *testing.T) {
	t.Parallel()

	_, err := NewServer(nil)
	if !errors.Is(err, ErrNilSource) {
		t.Fatalf("NewServer(nil) = %v, want ErrNilSource", err)
	}
}

func TestNewServerRejectsTypedNilSource(t *testing.T) {
	t.Parallel()

	var source *eval.Evaluator
	_, err := NewServer(source)
	if !errors.Is(err, ErrNilSource) {
		t.Fatalf("NewServer(typed nil) = %v, want ErrNilSource", err)
	}
}

func TestNewServerCreatesDefaultServiceMapping(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))
	if got := server.Services(); !sameStrings(got, []string{""}) {
		t.Fatalf("Services() = %v, want [empty]", got)
	}
	if !server.HasService("") {
		t.Fatal("HasService(empty) = false, want true")
	}
}

func TestNewServerRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []Option
		wantErr error
	}{
		{
			name: "duplicate service",
			opts: []Option{
				WithService("dup", health.TargetReady),
				WithService("dup", health.TargetLive),
			},
			wantErr: ErrDuplicateService,
		},
		{
			name:    "invalid watch interval",
			opts:    []Option{WithWatchInterval(0)},
			wantErr: ErrInvalidWatchInterval,
		},
		{
			name:    "invalid max list services",
			opts:    []Option{WithMaxListServices(0)},
			wantErr: ErrInvalidMaxListServices,
		},
		{
			name:    "nil option",
			opts:    []Option{nil},
			wantErr: ErrNilOption,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewServer(healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)), tc.opts...)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewServer() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestServerServicesReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	server := mustNewServer(
		t,
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
		WithService("custom", health.TargetLive),
	)
	services := server.Services()
	services[0] = "mutated"

	if got := server.Services(); !sameStrings(got, []string{"", "custom"}) {
		t.Fatalf("Services() after mutation = %v, want [ custom]", got)
	}
}

func TestServerHasServiceAndTarget(t *testing.T) {
	t.Parallel()

	server := mustNewServer(
		t,
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
		WithService("live", health.TargetLive),
	)

	if !server.HasService("live") {
		t.Fatal("HasService(live) = false, want true")
	}
	if server.HasService("missing") {
		t.Fatal("HasService(missing) = true, want false")
	}

	target, ok := server.Target("live")
	if !ok || target != health.TargetLive {
		t.Fatalf("Target(live) = %s, %v; want live true", target, ok)
	}
	target, ok = server.Target("missing")
	if ok || target != health.TargetUnknown {
		t.Fatalf("Target(missing) = %s, %v; want unknown false", target, ok)
	}
}

func TestNilServerReadMethods(t *testing.T) {
	t.Parallel()

	var server *Server
	if services := server.Services(); services != nil {
		t.Fatalf("Services() = %v, want nil", services)
	}
	if server.HasService("") {
		t.Fatal("HasService(empty) = true, want false")
	}
	target, ok := server.Target("")
	if ok || target != health.TargetUnknown {
		t.Fatalf("Target(empty) = %s, %v; want unknown false", target, ok)
	}
}
