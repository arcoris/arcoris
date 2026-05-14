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

package eval

import "fmt"

// PanicError preserves checker panic details as an internal health.Result.Cause.
//
// PanicError is intended for owner-controlled diagnostics and tests. Public
// adapters MUST NOT expose PanicError by default because panic values and stack
// traces can contain implementation details.
type PanicError struct {
	// Value is the recovered panic value.
	Value any

	// Stack is the stack captured when the panic was recovered.
	Stack []byte
}

// Error returns the panic classification message.
//
// Error intentionally does not include Stack. Callers that have permission to
// inspect internal diagnostics can access Stack through errors.As.
func (e PanicError) Error() string {
	return fmt.Sprintf("health check panicked: %v", e.Value)
}
