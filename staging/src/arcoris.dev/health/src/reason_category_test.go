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

import "testing"

func TestReasonCategoryPredicatesDoNotClassifyCustomReasons(t *testing.T) {
	t.Parallel()

	custom := Reason("custom_reason")
	predicates := []func(Reason) bool{
		Reason.IsObservationReason,
		Reason.IsExecutionReason,
		Reason.IsLifecycleReason,
		Reason.IsDependencyReason,
		Reason.IsControlReason,
		Reason.IsFreshnessReason,
		Reason.IsConnectivityReason,
		Reason.IsConfigurationReason,
	}

	for index, predicate := range predicates {
		if predicate(custom) {
			t.Fatalf("predicate %d classified custom reason", index)
		}
	}
}
