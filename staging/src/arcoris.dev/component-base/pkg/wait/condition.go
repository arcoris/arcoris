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

package wait

import "context"

// ConditionFunc evaluates whether a wait condition has been satisfied.
//
// A condition returns done=true when the caller should stop waiting
// successfully. It returns done=false with a nil error when the caller should
// keep waiting. It returns a non-nil error when waiting must stop and the error
// must be returned to the wait owner.
//
// The context is owned by the wait operation. Condition implementations SHOULD
// observe ctx when evaluation can block, perform I/O, acquire resources, or
// depend on cancellation or deadlines. Pure in-memory conditions MAY ignore ctx.
//
// ConditionFunc implementations MUST be safe to call repeatedly by the wait loop
// that owns them. They MUST NOT assume an exact number of evaluations because a
// wait may stop early due to success, error, context cancellation, timeout, or
// shutdown.
//
// Callers MUST pass a non-nil context. Use context.Background when no specific
// cancellation scope is available.
type ConditionFunc func(ctx context.Context) (done bool, err error)
