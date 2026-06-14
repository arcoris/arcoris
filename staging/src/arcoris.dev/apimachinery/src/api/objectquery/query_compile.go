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

// Compile validates and canonicalizes query into a predicate.
//
// The returned Predicate contains only validated canonical selector state. A
// zero Query compiles successfully and matches every item.
func Compile(query Query) (Predicate, error) {
	if err := query.Identity.Validate(); err != nil {
		return Predicate{}, err
	}

	labels, err := query.Labels.canonical()
	if err != nil {
		return Predicate{}, wrapf(
			"query.labels",
			ErrInvalidQuery,
			ErrorReasonInvalidQuery,
			err,
			"label selector is invalid",
		)
	}

	annotations, err := query.Annotations.canonical()
	if err != nil {
		return Predicate{}, wrapf(
			"query.annotations",
			ErrInvalidQuery,
			ErrorReasonInvalidQuery,
			err,
			"annotation selector is invalid",
		)
	}

	return Predicate{
		identity:    query.Identity,
		labels:      labels,
		annotations: annotations,
	}, nil
}
