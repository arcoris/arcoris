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
	"strings"
	"testing"

	"arcoris.dev/component-base/pkg/health"
	"arcoris.dev/component-base/pkg/healthtest"
)

func TestDefaultServiceMapping(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))
	target, ok := server.Target("")
	if !ok {
		t.Fatal("Target(empty) ok = false, want true")
	}
	if target != health.TargetReady {
		t.Fatalf("Target(empty) = %s, want ready", target)
	}
}

func TestTargetServiceMappings(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)), WithTargetServices())
	wantServices := []string{"", "startup", "live", "ready"}

	if got := server.Services(); !sameStrings(got, wantServices) {
		t.Fatalf("Services() = %v, want %v", got, wantServices)
	}

	tests := map[string]health.Target{
		"startup": health.TargetStartup,
		"live":    health.TargetLive,
		"ready":   health.TargetReady,
	}
	for service, want := range tests {
		target, ok := server.Target(service)
		if !ok || target != want {
			t.Fatalf("Target(%q) = %s, %v; want %s true", service, target, ok, want)
		}
	}
}

func TestServiceNameValidationAcceptsEmptyAndDottedNames(t *testing.T) {
	t.Parallel()

	mappings := []ServiceMapping{
		{Service: "", Target: health.TargetReady, Policy: health.ReadyPolicy()},
		{
			Service: "arcoris.control.v1.ControlPlane",
			Target:  health.TargetLive,
			Policy:  health.LivePolicy(),
		},
	}

	if _, err := normalizeServiceMappings(mappings); err != nil {
		t.Fatalf("normalizeServiceMappings() = %v, want nil", err)
	}
}

func TestServiceNameValidationRejectsInvalidNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		service string
	}{
		{"leading whitespace", " service"},
		{"trailing whitespace", "service "},
		{"control character", "service\nname"},
		{"too long", strings.Repeat("a", maxServiceNameLength+1)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewServer(
				healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
				WithService(tc.service, health.TargetReady),
			)
			if !errors.Is(err, ErrInvalidService) {
				t.Fatalf("NewServer() = %v, want ErrInvalidService", err)
			}
		})
	}
}

func TestServiceMappingRejectsInvalidTarget(t *testing.T) {
	t.Parallel()

	_, err := NewServer(
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
		WithService("invalid", health.TargetUnknown),
	)
	if !errors.Is(err, health.ErrInvalidTarget) {
		t.Fatalf("NewServer() = %v, want health.ErrInvalidTarget", err)
	}
}

func TestServiceMappingRejectsDuplicateService(t *testing.T) {
	t.Parallel()

	_, err := NewServer(
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
		WithService("duplicate", health.TargetReady),
		WithService("duplicate", health.TargetLive),
	)
	if !errors.Is(err, ErrDuplicateService) {
		t.Fatalf("NewServer() = %v, want ErrDuplicateService", err)
	}

	var duplicateErr DuplicateServiceError
	if !errors.As(err, &duplicateErr) {
		t.Fatalf("errors.As(%T, DuplicateServiceError) = false, want true", err)
	}
	if duplicateErr.Service != "duplicate" || duplicateErr.Index != 2 || duplicateErr.PreviousIndex != 1 {
		t.Fatalf("DuplicateServiceError = %+v, want service duplicate indexes 2 and 1", duplicateErr)
	}
}

func TestServiceIndexPreservesOrder(t *testing.T) {
	t.Parallel()

	server := mustNewServer(
		t,
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
		WithService("alpha", health.TargetReady),
		WithService("beta", health.TargetLive),
		WithService("gamma", health.TargetStartup),
	)

	want := []string{"", "alpha", "beta", "gamma"}
	if got := server.Services(); !sameStrings(got, want) {
		t.Fatalf("Services() = %v, want %v", got, want)
	}
}

// sameStrings reports whether two string slices contain the same ordered values.
func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}
