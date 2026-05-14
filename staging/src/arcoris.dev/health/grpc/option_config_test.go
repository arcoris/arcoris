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
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	config := defaultConfig()
	if len(config.services) != 1 || config.services[0].Service != "" {
		t.Fatalf("services = %+v, want default service only", config.services)
	}
	if config.watchInterval != defaultWatchInterval {
		t.Fatalf("watchInterval = %s, want %s", config.watchInterval, defaultWatchInterval)
	}
	if nilClock(config.clock) {
		t.Fatal("clock is nil")
	}
	if config.maxListServices != defaultMaxListServices {
		t.Fatalf("maxListServices = %d, want %d", config.maxListServices, defaultMaxListServices)
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*config)
		wantErr error
	}{
		{
			name: "valid",
		},
		{
			name: "nil clock",
			mutate: func(config *config) {
				config.clock = nil
			},
			wantErr: ErrNilClock,
		},
		{
			name: "invalid watch interval",
			mutate: func(config *config) {
				config.watchInterval = 0
			},
			wantErr: ErrInvalidWatchInterval,
		},
		{
			name: "invalid max list services",
			mutate: func(config *config) {
				config.maxListServices = 0
			},
			wantErr: ErrInvalidMaxListServices,
		},
		{
			name: "duplicate service",
			mutate: func(config *config) {
				config.services = append(config.services, ServiceMapping{
					Service: "",
					Target:  health.TargetReady,
					Policy:  health.ReadyPolicy(),
				})
			},
			wantErr: ErrDuplicateService,
		},
		{
			name: "invalid service mapping",
			mutate: func(config *config) {
				config.services = []ServiceMapping{{
					Service: " invalid",
					Target:  health.TargetReady,
					Policy:  health.ReadyPolicy(),
				}}
			},
			wantErr: ErrInvalidService,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			config := defaultConfig()
			if tc.mutate != nil {
				tc.mutate(&config)
			}

			err := config.validate()
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("validate() = %v, want nil", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("validate() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidateHelpers(t *testing.T) {
	t.Parallel()

	if err := validateWatchInterval(time.Second); err != nil {
		t.Fatalf("validateWatchInterval() = %v, want nil", err)
	}
	if err := validateMaxListServices(1); err != nil {
		t.Fatalf("validateMaxListServices() = %v, want nil", err)
	}
}
