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

// maxServiceNameLength bounds service names accepted by adapter configuration.
//
// gRPC service names are caller-chosen transport identifiers and can be longer
// than package-health identifiers. The bound prevents accidental unbounded
// configuration growth while staying permissive for dotted protobuf service
// names and schema-like strings.
const maxServiceNameLength = 512

// ServiceMapping binds one gRPC service name to a package-health target.
//
// Service is a transport identity. Target selects the health evaluation scope.
// Policy controls how that target's health.Status is converted to a gRPC
// serving status.
type ServiceMapping struct {
	// Service is the exact grpc.health.v1 service name.
	//
	// The empty string is the standard whole-server health service. Non-empty
	// names must be pre-trimmed and free of ASCII control characters. The
	// package intentionally does not impose lower_snake_case, GVK, GVR, or
	// package-health identifier syntax on transport names.
	Service string

	// Target is the concrete package-health target evaluated for Service.
	Target health.Target

	// Policy controls how Target status is interpreted for gRPC serving state.
	Policy health.TargetPolicy
}
