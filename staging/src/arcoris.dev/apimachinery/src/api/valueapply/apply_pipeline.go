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

package valueapply

// apply runs request validation, value validation, preparation, conflict
// handling, merge, and ownership update in that order.
func (a *applier) apply(req Request) (Result, error) {
	if err := validateRequestShape(req); err != nil {
		return Result{}, err
	}
	if err := a.validateValues(req); err != nil {
		return Result{}, err
	}

	prepared, err := a.prepare(req)
	if err != nil {
		return prepared.Result(), err
	}
	if err := a.rejectConflicts(req, prepared); err != nil {
		return prepared.Result(), err
	}
	if err := a.rejectUnsupportedForceTakeover(req, prepared); err != nil {
		return prepared.Result(), err
	}

	planned, err := a.planMergeFields(req, prepared)
	if err != nil {
		return prepared.Result(), err
	}

	merged, err := a.merge(req, planned)
	if err != nil {
		return planned.Result(), err
	}

	completed, err := a.updateOwnership(req, merged)
	if err != nil {
		return merged.Result(), err
	}

	return completed.Result(), nil
}
