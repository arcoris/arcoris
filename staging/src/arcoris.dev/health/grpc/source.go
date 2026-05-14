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
	"context"

	"arcoris.dev/health"
)

// Source evaluates one health target and returns a package-health report.
//
// Source is the only health model boundary used by healthgrpc. *health.Evaluator
// satisfies this interface. Implementations own check execution and report
// construction; healthgrpc owns only the gRPC adapter mapping.
type Source interface {
	// Evaluate synchronously evaluates target and returns the resulting report.
	Evaluate(ctx context.Context, target health.Target) (health.Report, error)
}
