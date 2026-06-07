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
	"context"
	"errors"

	"arcoris.dev/health"
)

// interruptedResult converts context interruption into an unknown health result.
func interruptedResult(name string, ctx context.Context) health.Result {
	err := ctx.Err()
	cause := context.Cause(ctx)
	if cause == nil {
		cause = err
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return health.Unknown(
			name,
			health.ReasonTimeout,
			"health check timed out",
		).WithCause(cause)
	}

	return health.Unknown(
		name,
		health.ReasonCanceled,
		"health check canceled",
	).WithCause(cause)
}
