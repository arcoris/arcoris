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

package objectquery

import "arcoris.dev/apimachinery/api/meta/labels"

// LabelRequirement is one canonical label selector requirement.
type LabelRequirement struct {
	// req stores the shared metadata requirement representation.
	req metadataRequirement
}

// LabelExists requires key to be present in object labels.
func LabelExists(key string) (LabelRequirement, error) {
	return newLabelRequirement(key, OperatorExists)
}

// LabelDoesNotExist requires key to be absent from object labels.
func LabelDoesNotExist(key string) (LabelRequirement, error) {
	return newLabelRequirement(key, OperatorDoesNotExist)
}

// LabelEquals requires key to have value.
func LabelEquals(key string, value string) (LabelRequirement, error) {
	return newLabelRequirement(key, OperatorEquals, value)
}

// LabelNotEquals requires key to be absent or to differ from value.
func LabelNotEquals(key string, value string) (LabelRequirement, error) {
	return newLabelRequirement(key, OperatorNotEquals, value)
}

// LabelIn requires key to have one of values.
func LabelIn(key string, values ...string) (LabelRequirement, error) {
	return newLabelRequirement(key, OperatorIn, values...)
}

// LabelNotIn requires key to be absent or outside values.
func LabelNotIn(key string, values ...string) (LabelRequirement, error) {
	return newLabelRequirement(key, OperatorNotIn, values...)
}

// newLabelRequirement constructs a validated label requirement.
func newLabelRequirement(key string, op Operator, values ...string) (LabelRequirement, error) {
	req, err := metadataRequirementFrom(
		"query.labels.requirement",
		key,
		op,
		values,
		validateLabelKey,
		validateLabelValue,
	)
	if err != nil {
		return LabelRequirement{}, err
	}

	return LabelRequirement{req: req}, nil
}

// validate checks label requirement structure.
func (r LabelRequirement) validate(path string) error {
	return r.req.validate(path, validateLabelKey, validateLabelValue)
}

// clone returns a detached requirement.
func (r LabelRequirement) clone() LabelRequirement {
	return LabelRequirement{req: r.req.clone()}
}

// validateLabelKey delegates label key syntax to api/meta/labels.
func validateLabelKey(key string) error {
	return labels.Key(key).ValidateLexical()
}

// validateLabelValue delegates label value syntax to api/meta/labels.
func validateLabelValue(value string) error {
	return labels.Value(value).ValidateLexical()
}
