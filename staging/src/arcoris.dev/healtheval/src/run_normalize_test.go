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
	"errors"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestNormalizeEvaluatedResult(t *testing.T) {
	t.Parallel()

	res := normalizeEvaluatedResult(
		health.Result{Status: health.StatusHealthy, Duration: -time.Second},
		"storage",
		testObserved,
		time.Second,
	)
	if res.Name != "storage" {
		t.Fatalf("name = %q, want storage", res.Name)
	}
	if res.Observed != testObserved {
		t.Fatalf("observed = %v, want %v", res.Observed, testObserved)
	}
	if res.Duration != time.Second {
		t.Fatalf("duration = %s, want 1s", res.Duration)
	}

	res = normalizeEvaluatedResult(
		health.Result{Status: health.StatusHealthy, Duration: -time.Second},
		"storage",
		testObserved,
		-time.Second,
	)
	if res.Duration != 0 {
		t.Fatalf("negative fallback duration = %s, want 0", res.Duration)
	}
}

func TestNormalizeEvaluatedResultRejectsMismatchedResultName(t *testing.T) {
	t.Parallel()

	res := normalizeEvaluatedResult(
		health.Healthy("database"),
		"storage",
		testObserved,
		time.Second,
	)

	if res.Name != "storage" || res.Status != health.StatusUnknown || res.Reason != health.ReasonMisconfigured {
		t.Fatalf("mismatched result normalization = %+v, want unknown misconfigured storage", res)
	}
	if !errors.Is(res.Cause, ErrMismatchedCheckResult) {
		t.Fatalf("cause = %v, want ErrMismatchedCheckResult", res.Cause)
	}
}

func TestNormalizeEvaluatedResultRejectsInvalidReason(t *testing.T) {
	t.Parallel()

	res := normalizeEvaluatedResult(
		health.Result{Status: health.StatusHealthy, Reason: health.Reason("bad-reason")},
		"storage",
		testObserved,
		time.Second,
	)

	if res.Name != "storage" || res.Status != health.StatusUnknown || res.Reason != health.ReasonMisconfigured {
		t.Fatalf("invalid reason normalization = %+v, want unknown misconfigured storage", res)
	}
	if !errors.Is(res.Cause, ErrInvalidCheckResult) {
		t.Fatalf("cause = %v, want ErrInvalidCheckResult", res.Cause)
	}
}
