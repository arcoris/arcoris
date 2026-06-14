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

// AnnotationExists requires key to be present in object annotations.
func AnnotationExists(key string) (AnnotationRequirement, error) {
	return newAnnotationRequirement(key, OperatorExists)
}

// AnnotationDoesNotExist requires key to be absent from object annotations.
func AnnotationDoesNotExist(key string) (AnnotationRequirement, error) {
	return newAnnotationRequirement(key, OperatorDoesNotExist)
}

// AnnotationEquals requires key to have value.
func AnnotationEquals(key string, value string) (AnnotationRequirement, error) {
	return newAnnotationRequirement(key, OperatorEquals, value)
}

// AnnotationNotEquals requires key to be absent or to differ from value.
func AnnotationNotEquals(key string, value string) (AnnotationRequirement, error) {
	return newAnnotationRequirement(key, OperatorNotEquals, value)
}

// AnnotationIn requires key to have one of values.
func AnnotationIn(key string, values ...string) (AnnotationRequirement, error) {
	return newAnnotationRequirement(key, OperatorIn, values...)
}

// AnnotationNotIn requires key to be absent or outside values.
func AnnotationNotIn(key string, values ...string) (AnnotationRequirement, error) {
	return newAnnotationRequirement(key, OperatorNotIn, values...)
}

// newAnnotationRequirement constructs a validated annotation requirement.
func newAnnotationRequirement(key string, op Operator, values ...string) (AnnotationRequirement, error) {
	req, err := metadataRequirementFrom(
		"query.annotations.requirement",
		key,
		op,
		values,
		validateAnnotationKey,
		validateAnnotationValue,
	)
	if err != nil {
		return AnnotationRequirement{}, err
	}

	return AnnotationRequirement{req: req}, nil
}
