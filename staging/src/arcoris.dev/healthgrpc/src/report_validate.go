// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package healthgrpc

import "arcoris.dev/health"

// validReportForTarget reports whether report is safe to map for target.
//
// health.Evaluator is an interface. Server therefore validates every evaluated
// report at the transport boundary before converting it to a gRPC serving
// status. Invalid, inconsistent, or target-mismatched reports are treated like
// evaluator failures and are never mapped as normal health states.
func validReportForTarget(report health.Report, target health.Target) bool {
	return report.Target == target &&
		report.IsValid() &&
		report.IsConsistent()
}
