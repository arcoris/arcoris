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

// Validator validates objects against a reusable object validation plan.
//
// Validator is immutable by convention. New stores the plan by value, and each
// Validate call runs an independent first-failure validation pipeline.
type Validator[D any, O any] struct {
	plan Plan[D, O]
}

// New returns a reusable object validator for plan.
//
// New does not validate plan eagerly. Static and conditional plan failures are
// reported by Validate so construction stays cheap and side-effect free.
func New[D any, O any](plan Plan[D, O]) Validator[D, O] {
	return Validator[D, O]{plan: plan}
}

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
	return New(plan).Validate(obj)
}

// Validate checks whether obj conforms to the validator's plan.
func (v Validator[D, O]) Validate(obj object.Object[D, O]) error {
	if err := validateStaticPlanShape(v.plan); err != nil {
		return err
	}

	if err := validateMetadata(obj); err != nil {
		return err
	}

	gvk := obj.GroupVersionKind()
	if err := validateResourceMatch(gvk, v.plan.Resource); err != nil {
		return err
	}

	version, err := resolveVersion(gvk, v.plan.Resource)
	if err != nil {
		return err
	}

	if err := validateScope(obj, v.plan.Resource); err != nil {
		return err
	}

	if err := validateDesired(obj.Desired, version, v.plan); err != nil {
		return err
	}

	observed, hasObserved := obj.ObservedValue()
	if err := validateObserved(observed, hasObserved, version, v.plan); err != nil {
		return err
	}

	return nil
}
