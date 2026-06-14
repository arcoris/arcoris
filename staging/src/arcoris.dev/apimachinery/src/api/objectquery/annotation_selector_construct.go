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

// NewAnnotationSelector constructs a canonical annotation selector.
//
// Exact duplicate requirements are removed. Requirements are sorted
// deterministically by key, operator, and values.
func NewAnnotationSelector(requirements ...AnnotationRequirement) (AnnotationSelector, error) {
	raw := make([]metadataRequirement, 0, len(requirements))
	for _, req := range requirements {
		raw = append(raw, req.req.clone())
	}
	canonical, err := canonicalMetadataRequirements(
		"query.annotations.requirements",
		raw,
		validateAnnotationKey,
		validateAnnotationValue,
	)
	if err != nil {
		return AnnotationSelector{}, wrapf(
			"query.annotations",
			ErrInvalidSelector,
			ErrorReasonInvalidSelector,
			err,
			"annotation selector is invalid",
		)
	}

	return AnnotationSelector{requirements: annotationRequirementsFromMetadata(canonical)}, nil
}
