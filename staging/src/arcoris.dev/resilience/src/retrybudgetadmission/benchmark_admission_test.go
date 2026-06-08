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

package retrybudgetadmission

import (
	"testing"

	"arcoris.dev/admission"
	"arcoris.dev/resilience/retrybudget"
)

var benchmarkResult admission.Result[admission.NoGrant, retrybudget.Decision]

func BenchmarkAdmissionResultAllowed(b *testing.B) {
	decision := validAllowedDecision()
	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = AdmissionResult(decision)
	}
}

func BenchmarkAdmissionResultDenied(b *testing.B) {
	decision := validDeniedDecision()
	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = AdmissionResult(decision)
	}
}

func BenchmarkAdmitterTryAdmitAllowed(b *testing.B) {
	admitter := New(alwaysDecisionBudget{decision: validAllowedDecision()})
	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = admitter.TryAdmit(Request{})
	}
}

func BenchmarkAdmitterTryAdmitDenied(b *testing.B) {
	admitter := New(alwaysDecisionBudget{decision: validDeniedDecision()})
	b.ReportAllocs()
	for b.Loop() {
		benchmarkResult = admitter.TryAdmit(Request{})
	}
}

type alwaysDecisionBudget struct {
	decision retrybudget.Decision
}

func (b alwaysDecisionBudget) TryAdmitRetry() retrybudget.Decision {
	return b.decision
}
