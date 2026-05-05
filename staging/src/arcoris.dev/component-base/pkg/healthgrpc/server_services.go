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

import "arcoris.dev/component-base/pkg/health"

// Services returns configured gRPC health service names in deterministic order.
//
// The returned slice is detached from the Server. Callers may sort, append to,
// or modify it without mutating the adapter's service index.
func (s *Server) Services() []string {
	if s == nil || len(s.order) == 0 {
		return nil
	}

	services := make([]string, len(s.order))
	copy(services, s.order)

	return services
}

// HasService reports whether service is configured on the adapter.
//
// Service names are matched exactly. healthgrpc does not trim, lowercase,
// canonicalize, or reinterpret names at read time because service names are
// transport identities chosen by configuration.
func (s *Server) HasService(service string) bool {
	_, ok := s.service(service)
	return ok
}

// Target returns the package-health target configured for service.
//
// The boolean is false when the service name is not configured or when the
// method is called on a nil Server. Unknown services intentionally do not fall
// back to the default service; the gRPC health protocol treats them as distinct
// transport identities.
func (s *Server) Target(service string) (health.Target, bool) {
	mapping, ok := s.service(service)
	if !ok {
		return health.TargetUnknown, false
	}

	return mapping.Target, true
}

// service returns the immutable mapping for service.
//
// It is the shared lookup boundary for public readers and RPC handlers, keeping
// nil-Server behavior stable and avoiding direct map access outside this file.
func (s *Server) service(service string) (ServiceMapping, bool) {
	if s == nil {
		return ServiceMapping{}, false
	}

	mapping, ok := s.services[service]
	return mapping, ok
}
