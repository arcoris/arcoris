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

	config := defaultConfig()
	err := applyOptions(&config, nil)
	if !errors.Is(err, ErrNilOption) {
		t.Fatalf("applyOptions(nil) = %v, want ErrNilOption", err)
	}
}

func TestServiceOptions(t *testing.T) {
	t.Parallel()

	config := defaultConfig()
	err := applyOptions(
		&config,
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

	if len(config.services) != 4 {
		t.Fatalf("services length = %d, want 4", len(config.services))
	}
	if config.services[1].Service != "live-service" || config.services[1].Policy != health.LivePolicy() {
		t.Fatalf("WithService mapping = %+v, want live default policy", config.services[1])
	}
	if !config.services[2].Policy.AllowDegraded {
		t.Fatalf("WithServicePolicy policy = %+v, want AllowDegraded", config.services[2].Policy)
	}
	if config.services[3].Target != health.TargetStartup {
		t.Fatalf("WithServices target = %s, want startup", config.services[3].Target)
	}
}

func TestDefaultServiceOptions(t *testing.T) {
	t.Parallel()

	config := defaultConfig()
	err := applyOptions(
		&config,
		WithService("custom", health.TargetReady),
		WithDefaultService(health.TargetLive),
		WithDefaultServicePolicy(health.TargetReady, health.ReadyPolicy().WithDegraded(true)),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if config.services[0].Service != "" || config.services[0].Target != health.TargetReady {
		t.Fatalf("default service = %+v, want ready", config.services[0])
	}
	if !config.services[0].Policy.AllowDegraded {
		t.Fatalf("default policy = %+v, want AllowDegraded", config.services[0].Policy)
	}
	if config.services[1].Service != "custom" {
		t.Fatalf("second service = %q, want custom", config.services[1].Service)
	}
}

func TestWithTargetServicesOption(t *testing.T) {
	t.Parallel()

	config := defaultConfig()
	if err := WithTargetServices()(&config); err != nil {
		t.Fatalf("WithTargetServices() = %v, want nil", err)
	}

	want := []string{"", "startup", "live", "ready"}
	var got []string
	for _, mapping := range config.services {
		got = append(got, mapping.Service)
	}
	if !sameStrings(got, want) {
		t.Fatalf("services = %v, want %v", got, want)
	}
}

func TestScheduleOptions(t *testing.T) {
	t.Parallel()

	config := defaultConfig()
	err := applyOptions(
		&config,
		WithWatchInterval(time.Second),
		WithWatchInterval(2*time.Second),
		WithMaxListServices(10),
		WithMaxListServices(20),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}
	if config.watchInterval != 2*time.Second {
		t.Fatalf("watchInterval = %s, want 2s", config.watchInterval)
	}
	if config.maxListServices != 20 {
		t.Fatalf("maxListServices = %d, want 20", config.maxListServices)
	}
}

func TestOptionsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		option  Option
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

			config := defaultConfig()
			err := tc.option(&config)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("option() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestReplaceDefaultServiceMappingAppendsWhenMissing(t *testing.T) {
	t.Parallel()

	config := config{}
	replaceDefaultServiceMapping(&config, ServiceMapping{
		Service: "",
		Target:  health.TargetReady,
		Policy:  health.ReadyPolicy(),
	})

	if len(config.services) != 1 || config.services[0].Target != health.TargetReady {
		t.Fatalf("services = %+v, want appended ready default", config.services)
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
