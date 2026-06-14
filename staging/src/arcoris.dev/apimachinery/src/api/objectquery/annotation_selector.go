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

import (
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/objectstore"
)

// AnnotationSelector is a canonical AND-set of annotation requirements.
type AnnotationSelector struct {
	// requirements are sorted by key, operator, and values.
	requirements []AnnotationRequirement
}

// NewAnnotationSelector constructs a canonical annotation selector.
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

// IsZero reports whether s has no annotation requirements.
func (s AnnotationSelector) IsZero() bool {
	return len(s.requirements) == 0
}

// Validate checks selector structure.
func (s AnnotationSelector) Validate() error {
	_, err := NewAnnotationSelector(s.requirements...)
	return err
}

// canonical returns a validated canonical copy of s.
func (s AnnotationSelector) canonical() (AnnotationSelector, error) {
	return NewAnnotationSelector(s.requirements...)
}

// match checks item annotations against the selector.
func (s AnnotationSelector) match(item objectstore.ListItem) bool {
	return matchAnnotationRequirements(s.requirements, func(key string) (string, bool) {
		value, ok := item.State.Object.ObjectMeta.Annotations.Get(annotations.Key(key))
		return value.String(), ok
	})
}

// annotationRequirementsFromMetadata adapts shared metadata requirements to annotations.
func annotationRequirementsFromMetadata(requirements []metadataRequirement) []AnnotationRequirement {
	if len(requirements) == 0 {
		return nil
	}

	out := make([]AnnotationRequirement, len(requirements))
	for i, req := range requirements {
		out[i] = AnnotationRequirement{req: req.clone()}
	}

	return out
}

// matchAnnotationRequirements applies AND semantics without allocating per item.
func matchAnnotationRequirements(requirements []AnnotationRequirement, lookup func(string) (string, bool)) bool {
	for _, req := range requirements {
		if !req.req.match(lookup) {
			return false
		}
	}

	return true
}
