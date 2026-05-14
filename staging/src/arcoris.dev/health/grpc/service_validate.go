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
	"strings"

	"arcoris.dev/health"
)

// normalizeServiceMappings validates mappings and returns a detached copy.
//
// Validation is centralized here so options can fail early while NewServer
// still validates the final ordered set after all options have been applied.
// Duplicate detection happens only on the final set because options are
// intentionally order-sensitive and may replace the default mapping.
func normalizeServiceMappings(mappings []ServiceMapping) ([]ServiceMapping, error) {
	normalized := make([]ServiceMapping, 0, len(mappings))
	seen := make(map[string]int, len(mappings))

	for index, mapping := range mappings {
		if err := validateServiceMapping(mapping, index); err != nil {
			return nil, err
		}
		if previous, exists := seen[mapping.Service]; exists {
			return nil, DuplicateServiceError{
				Service:       mapping.Service,
				Index:         index,
				PreviousIndex: previous,
			}
		}

		seen[mapping.Service] = index
		normalized = append(normalized, mapping)
	}

	return normalized, nil
}

// validateServiceMapping validates one mapping without checking duplicates.
//
// The target check deliberately reuses package-health target errors. healthgrpc
// owns service-name validation, but health remains the owner of target
// concreteness and target error classification.
func validateServiceMapping(mapping ServiceMapping, index int) error {
	if err := validateServiceName(mapping.Service, index); err != nil {
		return err
	}
	if !mapping.Target.IsConcrete() {
		return health.InvalidTargetError{Target: mapping.Target}
	}

	return nil
}

// validateServiceName enforces only transport-safety rules for gRPC names.
//
// Empty is valid by protocol. Non-empty names must already be trimmed so config
// mistakes are visible, and ASCII control characters are rejected because they
// make diagnostics and transport metadata unsafe to inspect.
func validateServiceName(service string, index int) error {
	if service == "" {
		return nil
	}
	if service != strings.TrimSpace(service) {
		return InvalidServiceError{Service: service, Index: index, Reason: "service name must be trimmed"}
	}
	if len(service) > maxServiceNameLength {
		return InvalidServiceError{Service: service, Index: index, Reason: "service name is too long"}
	}
	for _, r := range service {
		if r < 0x20 || r == 0x7f {
			return InvalidServiceError{
				Service: service,
				Index:   index,
				Reason:  "service name contains an ASCII control character",
			}
		}
	}

	return nil
}
