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
	"time"

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
)

func TestApplyOptionsRejectsNilOption(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := applyOptions(&cfg, nil)
	if !errors.Is(err, ErrNilOption) {
		t.Fatalf("applyOptions(nil) = %v, want ErrNilOption", err)
	}
}

func TestServiceOptions(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := applyOptions(
		&cfg,
		WithService("live-service", health.TargetLive),
		WithServicePolicy("ready-degraded", health.TargetReady, health.ReadyPolicy().WithDegraded(true)),
		WithServices(ServiceMapping{
			Service: "startup-service",
			Target:  health.TargetStartup,
			Policy:  health.StartupPolicy(),
		}),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if len(cfg.services) != 4 {
		t.Fatalf("services length = %d, want 4", len(cfg.services))
	}
	if cfg.services[1].Service != "live-service" || cfg.services[1].Policy != health.LivePolicy() {
		t.Fatalf("WithService mapping = %+v, want live default policy", cfg.services[1])
	}
	if !cfg.services[2].Policy.AllowDegraded {
		t.Fatalf("WithServicePolicy policy = %+v, want AllowDegraded", cfg.services[2].Policy)
	}
	if cfg.services[3].Target != health.TargetStartup {
		t.Fatalf("WithServices target = %s, want startup", cfg.services[3].Target)
	}
}

func TestDefaultServiceOptions(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := applyOptions(
		&cfg,
		WithService("custom", health.TargetReady),
		WithDefaultService(health.TargetLive),
		WithDefaultServicePolicy(health.TargetReady, health.ReadyPolicy().WithDegraded(true)),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if cfg.services[0].Service != "" || cfg.services[0].Target != health.TargetReady {
		t.Fatalf("default service = %+v, want ready", cfg.services[0])
	}
	if !cfg.services[0].Policy.AllowDegraded {
		t.Fatalf("default policy = %+v, want AllowDegraded", cfg.services[0].Policy)
	}
	if cfg.services[1].Service != "custom" {
		t.Fatalf("second service = %q, want custom", cfg.services[1].Service)
	}
}

func TestWithTargetServicesOption(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	if err := WithTargetServices()(&cfg); err != nil {
		t.Fatalf("WithTargetServices() = %v, want nil", err)
	}

	want := []string{"", "startup", "live", "ready"}
	var got []string
	for _, mapping := range cfg.services {
		got = append(got, mapping.Service)
	}
	if !sameStrings(got, want) {
		t.Fatalf("services = %v, want %v", got, want)
	}
}

func TestScheduleOptions(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := applyOptions(
		&cfg,
		WithWatchInterval(time.Second),
		WithWatchInterval(2*time.Second),
		WithMaxListServices(10),
		WithMaxListServices(20),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}
	if cfg.watchInterval != 2*time.Second {
		t.Fatalf("watchInterval = %s, want 2s", cfg.watchInterval)
	}
	if cfg.maxListServices != 20 {
		t.Fatalf("maxListServices = %d, want 20", cfg.maxListServices)
	}
}

func TestOptionsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opt     Option
		wantErr error
	}{
		{"invalid service", WithService(" bad", health.TargetReady), ErrInvalidService},
		{"invalid target", WithService("bad-target", health.TargetUnknown), health.ErrInvalidTarget},
		{
			"invalid WithServices mapping",
			WithServices(ServiceMapping{Service: "bad", Target: health.TargetUnknown}),
			health.ErrInvalidTarget,
		},
		{
			"invalid default service target",
			WithDefaultServicePolicy(health.TargetUnknown, health.ReadyPolicy()),
			health.ErrInvalidTarget,
		},
		{"invalid watch interval", WithWatchInterval(0), ErrInvalidWatchInterval},
		{"invalid max list services", WithMaxListServices(0), ErrInvalidMaxListServices},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig()
			err := tc.opt(&cfg)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("option() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestReplaceDefaultServiceMappingAppendsWhenMissing(t *testing.T) {
	t.Parallel()

	cfg := config{}
	replaceDefaultServiceMapping(&cfg, ServiceMapping{
		Service: "",
		Target:  health.TargetReady,
		Policy:  health.ReadyPolicy(),
	})

	if len(cfg.services) != 1 || cfg.services[0].Target != health.TargetReady {
		t.Fatalf("services = %+v, want appended ready default", cfg.services)
	}
}

func TestOptionsAreAppliedInOrder(t *testing.T) {
	t.Parallel()

	server := mustNewServer(
		t,
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
		WithDefaultService(health.TargetLive),
		WithDefaultService(health.TargetReady),
		WithService("custom", health.TargetStartup),
		WithWatchInterval(time.Second),
		WithWatchInterval(2*time.Second),
	)

	target, ok := server.Target("")
	if !ok || target != health.TargetReady {
		t.Fatalf("Target(empty) = %s, %v; want ready true", target, ok)
	}
	if server.config.watchInterval != 2*time.Second {
		t.Fatalf("watchInterval = %s, want 2s", server.config.watchInterval)
	}
	if got := server.Services(); !sameStrings(got, []string{"", "custom"}) {
		t.Fatalf("Services() = %v, want [ custom]", got)
	}
}
