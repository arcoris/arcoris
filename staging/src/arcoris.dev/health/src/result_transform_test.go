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

package health

import (
	"errors"
	"testing"
	"time"
)

func TestResultTransformMethodsReturnCopies(t *testing.T) {
	t.Parallel()

	base := Result{Status: StatusHealthy}
	cause := errors.New("cause")
	changed := base.
		WithCause(cause).
		WithObserved(testObserved).
		WithDuration(time.Second).
		WithMessage("message").
		WithReason(ReasonFatal)

	if base.Cause != nil || !base.Observed.IsZero() || base.Duration != 0 || base.Message != "" || base.Reason != ReasonNone {
		t.Fatalf("base result mutated: %+v", base)
	}
	if changed.Cause != cause || changed.Observed != testObserved || changed.Duration != time.Second {
		t.Fatalf("changed result missing metadata: %+v", changed)
	}
}
