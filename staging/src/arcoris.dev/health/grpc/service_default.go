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

// defaultServiceMapping returns the standard whole-server gRPC health mapping.
//
// The gRPC health protocol reserves the empty service name for the whole
// server. healthgrpc maps that identity to readiness by default because ready
// is the least surprising target for callers asking whether the process should
// receive normal traffic.
func defaultServiceMapping() ServiceMapping {
	return ServiceMapping{
		Service: "",
		Target:  health.TargetReady,
		Policy:  health.ReadyPolicy(),
	}
}

// targetServiceMappings returns the opt-in built-in target service names.
//
// These mappings are not enabled by default so package owners can choose their
// public gRPC service namespace explicitly. When enabled, the transport names
// mirror the package-health targets without introducing a second target model.
func targetServiceMappings() []ServiceMapping {
	return []ServiceMapping{
		{Service: "startup", Target: health.TargetStartup, Policy: health.StartupPolicy()},
		{Service: "live", Target: health.TargetLive, Policy: health.LivePolicy()},
		{Service: "ready", Target: health.TargetReady, Policy: health.ReadyPolicy()},
	}
}
