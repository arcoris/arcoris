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

package health

import "context"

// Evaluator evaluates one health target and returns a target-level report.
//
// Evaluator is a root health contract. It is implemented by synchronous
// evaluators, cached evaluators, test fakes, adapters, and other report sources
// that can produce a Report for a concrete Target.
//
// Evaluator does not prescribe how checks are executed. Execution policy,
// timeouts, panic recovery, caching, scheduling, transport rendering, metrics,
// logging, lifecycle decisions, routing, admission, and restart policy belong to
// concrete implementations or higher-level packages.
type Evaluator interface {
	Evaluate(ctx context.Context, target Target) (Report, error)
}
