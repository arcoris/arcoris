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

package objectvalidation

import "arcoris.dev/apimachinery/api/object"

// Validate checks whether obj conforms to the already-resolved plan Resource.
//
// Validation is deterministic and returns the first failure: plan shape,
// metadata, resource group/kind match, version lookup, scope compatibility,
// desired surface validation, then observed surface validation. It does not
// mutate the object and does not perform request admission, serving, storage,
// conversion, defaulting, or catalog lookup.
func Validate[D any, O any](
	obj object.Object[D, O],
	plan Plan[D, O],
) error {
	if err := validatePlanShape(plan); err != nil {
		return err
	}

	if err := validateMetadata(obj); err != nil {
		return err
	}

	gvk := obj.GroupVersionKind()
	if err := validateResourceMatch(gvk, plan.Resource); err != nil {
		return err
	}

	version, err := resolveVersion(gvk, plan.Resource)
	if err != nil {
		return err
	}

	if err := validateScope(obj, plan.Resource); err != nil {
		return err
	}

	if err := validateDesired(obj.Desired, version, plan); err != nil {
		return err
	}

	if err := validateObserved(obj.Observed, version, plan); err != nil {
		return err
	}

	return nil
}
