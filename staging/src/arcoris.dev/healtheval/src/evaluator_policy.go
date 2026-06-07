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

package eval

import (
	"time"

	"arcoris.dev/health"
)

// timeoutFor returns the effective check timeout for target.
func (e *Evaluator) timeoutFor(target health.Target) time.Duration {
	if timeout, ok := e.targetTimeouts[target]; ok {
		return timeout
	}

	return e.defaultTimeout
}

// executionPolicyFor returns the effective check execution policy for target.
//
// health.Target-specific execution policy overrides the evaluator default. The
// returned policy is normalized at construction time.
func (e *Evaluator) executionPolicyFor(target health.Target) ExecutionPolicy {
	if policy, ok := e.targetExecutionPolicies[target]; ok {
		return policy
	}

	return e.executionPolicy
}
