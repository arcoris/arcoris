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

package objectreconciler

// Result is the outcome of one reconciliation attempt.
//
// Err nil means the attempt completed successfully. Err non-nil means the
// attempt failed. Result intentionally has no retry, requeue, priority, or
// terminal-failure fields; those policies belong to future queue and controller
// layers.
type Result struct {
	// Err is the reconciliation failure, or nil when the attempt completed
	// successfully.
	//
	// ReconcileOnce preserves source and user errors in Err without wrapping.
	Err error
}

// Success returns a successful reconciliation result with a nil Err field.
func Success() Result {
	return Result{}
}

// Failure returns a failed reconciliation result for err.
//
// Failure(nil) returns the same successful zero-error result as Success.
func Failure(err error) Result {
	return Result{Err: err}
}

// Failed reports whether r contains a non-nil Err.
func (r Result) Failed() bool {
	return r.Err != nil
}
