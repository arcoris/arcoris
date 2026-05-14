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

import "arcoris.dev/health"

// WithService appends a service mapping using package-health's default policy.
//
// The service name is a gRPC transport identity and the target is a concrete
// package-health target. The option validates the single mapping immediately;
// duplicate service names are checked after all options are applied so
// order-sensitive replacement options remain possible.
func WithService(service string, target health.Target) Option {
	return WithServicePolicy(service, target, health.DefaultPolicy(target))
}

// WithServicePolicy appends a service mapping with an explicit target policy.
//
// This option is useful when the transport surface should interpret a target
// differently from health.DefaultPolicy. It does not change the target's checks,
// evaluator execution policy, or package-health aggregation behavior.
func WithServicePolicy(service string, target health.Target, policy health.TargetPolicy) Option {
	return func(cfg *config) error {
		mapping := ServiceMapping{Service: service, Target: target, Policy: policy}
		if err := validateServiceMapping(mapping, len(cfg.services)); err != nil {
			return err
		}

		cfg.services = append(cfg.services, mapping)
		return nil
	}
}

// WithServices appends service mappings in the supplied order.
//
// The mappings are copied into adapter configuration by value. Callers may
// reuse or modify their input slice after NewServer returns without mutating the
// Server's immutable service index.
func WithServices(mappings ...ServiceMapping) Option {
	return func(cfg *config) error {
		for i, mapping := range mappings {
			if err := validateServiceMapping(mapping, len(cfg.services)+i); err != nil {
				return err
			}
		}

		cfg.services = append(cfg.services, mappings...)
		return nil
	}
}

// WithDefaultService replaces the standard whole-server service mapping.
//
// The empty gRPC service name remains the service identity; only the
// package-health target and default policy are changed.
func WithDefaultService(target health.Target) Option {
	return WithDefaultServicePolicy(target, health.DefaultPolicy(target))
}

// WithDefaultServicePolicy replaces the standard whole-server service mapping.
//
// The replacement preserves the original position of the first default mapping
// so service order remains deterministic. Additional default mappings created
// by earlier options are removed to keep the final config duplicate-free.
func WithDefaultServicePolicy(target health.Target, policy health.TargetPolicy) Option {
	return func(cfg *config) error {
		mapping := ServiceMapping{Service: "", Target: target, Policy: policy}
		if err := validateServiceMapping(mapping, 0); err != nil {
			return err
		}

		replaceDefaultServiceMapping(cfg, mapping)
		return nil
	}
}

// WithTargetServices appends startup, live, and ready service mappings.
//
// The option is explicit so the adapter does not publish target-specific names
// unless the owner wants that gRPC surface.
func WithTargetServices() Option {
	return func(cfg *config) error {
		cfg.services = append(cfg.services, targetServiceMappings()...)
		return nil
	}
}

// replaceDefaultServiceMapping replaces all empty-service mappings in cfg.
//
// The first empty-service slot keeps its position, which preserves caller order
// for List and Services. If the config has no empty-service mapping, the
// replacement is appended; this keeps the helper robust for direct config tests.
func replaceDefaultServiceMapping(cfg *config, mapping ServiceMapping) {
	services := cfg.services[:0]
	replaced := false

	for _, existing := range cfg.services {
		if existing.Service == "" {
			if !replaced {
				services = append(services, mapping)
				replaced = true
			}
			continue
		}

		services = append(services, existing)
	}
	if !replaced {
		services = append(services, mapping)
	}

	cfg.services = services
}
