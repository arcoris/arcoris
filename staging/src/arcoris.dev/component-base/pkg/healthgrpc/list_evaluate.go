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

	"arcoris.dev/component-base/pkg/health"
)

// targetEvaluation is the per-target List cache entry.
//
// List may expose several service names for the same health.Target. This value
// lets the method evaluate that target once, remember whether Source failed,
// and apply each service's own policy afterward.
type targetEvaluation struct {
	// status is the package-health status returned by Source.
	status health.Status

	// failed records Source failure without storing or exposing the raw error.
	failed bool
}

// evaluateTarget evaluates target for List and normalizes failures.
//
// A Source error is intentionally reduced to failed=true. List maps affected
// service responses to gRPC UNKNOWN and never exposes the raw error text.
func (s *Server) evaluateTarget(ctx context.Context, target health.Target) targetEvaluation {
	report, err := s.source.Evaluate(ctx, target)
	if err != nil {
		return targetEvaluation{failed: true}
	}

	return targetEvaluation{status: report.Status}
}
