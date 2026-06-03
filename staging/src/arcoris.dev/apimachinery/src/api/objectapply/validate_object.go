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

package objectapply

import "arcoris.dev/apimachinery/api/objectvalidation"

// validateObjectMeta checks envelope metadata before policy comparisons.
//
// object.ValidateMeta covers only TypeMeta/ObjectMeta structure. Desired and
// Observed validation stay in objectvalidation/valuevalidation.
func validateObjectMeta(path string, obj ValueObject, reason ErrorReason) error {
	if err := obj.ValidateMeta(); err != nil {
		return wrapAt(
			path,
			ErrInvalidObject,
			reason,
			"object metadata is invalid",
			err,
		)
	}

	return nil
}

// validateObject delegates resource conformance to api/objectvalidation.
//
// This keeps objectapply from becoming a second object validator. The local
// wrapper only maps validation failures into the objectapply error taxonomy.
func (a applier) validateObject(
	path string,
	obj ValueObject,
	reason ErrorReason,
	req Request,
) error {
	if err := objectvalidation.Validate(obj, a.validationPlan(req)); err != nil {
		return wrapAt(
			path,
			ErrInvalidObject,
			reason,
			"object does not satisfy resource contract",
			err,
		)
	}

	return nil
}
