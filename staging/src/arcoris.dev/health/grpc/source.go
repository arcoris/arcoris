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

// Source evaluates one health target and returns a package-health report.
//
// Source is an alias for health.Evaluator kept as the transport-facing name in
// healthgrpc documentation. Implementations own check execution and report
// construction; healthgrpc owns only the gRPC adapter mapping.
type Source = health.Evaluator
