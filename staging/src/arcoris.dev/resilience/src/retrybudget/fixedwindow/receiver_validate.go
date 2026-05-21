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


package fixedwindow

const (
	// nilLimiterPanic is the stable panic value for methods called on nil
	// *Limiter receivers.
	nilLimiterPanic = "retrybudget.fixedwindow.Limiter: nil limiter"

	// uninitializedLimiterPanic is the stable panic value for methods called on
	// zero or otherwise not-New-created Limiter values.
	uninitializedLimiterPanic = "retrybudget.fixedwindow.Limiter: uninitialized limiter"
)

// requireReady validates receiver ownership before public methods touch mutable
// limiter state.
//
// New validates policy configuration once. This guard is narrower: it only
// detects nil receivers and Limiter values that were not initialized by New.
func (l *Limiter) requireReady() {
	if l == nil {
		panic(nilLimiterPanic)
	}
	if l.published == nil || l.cfg.clock == nil {
		panic(uninitializedLimiterPanic)
	}
}
