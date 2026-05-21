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

// NewServer builds a standard gRPC health service adapter over source.
//
// Construction validates the full adapter configuration and builds an immutable
// service index. It does not start goroutines, allocate Watch streams, register
// itself with a grpc.Server, or ask source to evaluate anything. Evaluation is
// owned by Check, List, and Watch at request time.
func NewServer(source health.Evaluator, opts ...Option) (*Server, error) {
	if nilSource(source) {
		return nil, ErrNilSource
	}

	cfg := defaultConfig()
	if err := applyOptions(&cfg, opts...); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	mappings, err := normalizeServiceMappings(cfg.services)
	if err != nil {
		return nil, err
	}

	services, order := buildServiceIndex(mappings)
	cfg.services = mappings

	return &Server{
		source:   source,
		services: services,
		order:    order,
		config:   cfg,
	}, nil
}

// buildServiceIndex creates the private lookup structures used by Server.
//
// The input must already be validated for duplicate service names. Keeping the
// map and order slice separate lets request paths use O(1) lookup without
// losing deterministic configuration order for Services and List.
func buildServiceIndex(mappings []ServiceMapping) (map[string]ServiceMapping, []string) {
	services := make(map[string]ServiceMapping, len(mappings))
	order := make([]string, 0, len(mappings))

	for _, mapping := range mappings {
		services[mapping.Service] = mapping
		order = append(order, mapping.Service)
	}

	return services, order
}
