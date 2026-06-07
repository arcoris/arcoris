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

// AggregateStatus returns the aggregate status for results.
//
// Aggregation uses the root health severity ordering. Empty result slices return
// StatusUnknown because no affirmative health observation exists. Invalid result
// statuses dominate valid statuses so corrupted or caller-controlled input is
// not silently hidden by aggregation.
func AggregateStatus(results []Result) Status {
	if len(results) == 0 {
		return StatusUnknown
	}

	status := StatusHealthy
	for _, result := range results {
		if result.Status.MoreSevereThan(status) {
			status = result.Status
		}
	}

	return status
}

// AggregateStatus returns the status computed from r.Checks.
func (r Report) AggregateStatus() Status {
	return AggregateStatus(r.Checks)
}

// IsConsistent reports whether r.Status matches the aggregate status of r.Checks.
//
// Consistency is separate from structural validity. A caller-owned report can be
// structurally valid while carrying a stale or intentionally supplied aggregate
// status. Empty concrete reports are consistent when their status is
// StatusUnknown.
func (r Report) IsConsistent() bool {
	if !r.IsValid() {
		return false
	}
	if r.Target == TargetUnknown {
		return true
	}

	return r.Status == r.AggregateStatus()
}

// WithAggregateStatus returns a copy of r with Status recomputed from Checks.
func (r Report) WithAggregateStatus() Report {
	r.Status = r.AggregateStatus()
	return r
}
