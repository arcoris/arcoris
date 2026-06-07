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

package healthgate

import "arcoris.dev/health"

// normalizeGateResult validates result for storage in a gate and fills an empty
// result name with the gate name.
func normalizeGateResult(name string, result health.Result) (health.Result, error) {
	if result.Name == "" {
		result.Name = name
	} else if result.Name != name {
		return health.Result{}, MismatchedGateResultError{
			GateName:   name,
			ResultName: result.Name,
		}
	}

	if !result.IsValid() {
		return health.Result{}, InvalidGateResultError{
			GateName: name,
			Result:   result,
		}
	}

	return result, nil
}
