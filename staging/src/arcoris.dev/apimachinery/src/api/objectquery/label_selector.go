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
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

// LabelSelector is a canonical AND-set of label requirements.
type LabelSelector struct {
	// requirements are sorted by key, operator, and values.
	requirements []LabelRequirement
}

// NewLabelSelector constructs a canonical label selector.
func NewLabelSelector(requirements ...LabelRequirement) (LabelSelector, error) {
	raw := make([]metadataRequirement, 0, len(requirements))
	for _, req := range requirements {
		raw = append(raw, req.req.clone())
	}
	canonical, err := canonicalMetadataRequirements(
		"query.labels.requirements",
		raw,
		validateLabelKey,
		validateLabelValue,
	)
	if err != nil {
		return LabelSelector{}, wrapf(
			"query.labels",
			ErrInvalidSelector,
			ErrorReasonInvalidSelector,
			err,
			"label selector is invalid",
		)
	}

	return LabelSelector{requirements: labelRequirementsFromMetadata(canonical)}, nil
}

// IsZero reports whether s has no label requirements.
func (s LabelSelector) IsZero() bool {
	return len(s.requirements) == 0
}

// Validate checks selector structure.
func (s LabelSelector) Validate() error {
	_, err := NewLabelSelector(s.requirements...)
	return err
}

// canonical returns a validated canonical copy of s.
func (s LabelSelector) canonical() (LabelSelector, error) {
	return NewLabelSelector(s.requirements...)
}

// match checks item labels against the selector.
func (s LabelSelector) match(item objectstore.ListItem) bool {
	return matchLabelRequirements(s.requirements, func(key string) (string, bool) {
		value, ok := item.State.Object.ObjectMeta.Labels.Get(labels.Key(key))
		return value.String(), ok
	})
}

// labelRequirementsFromMetadata adapts shared metadata requirements to labels.
func labelRequirementsFromMetadata(requirements []metadataRequirement) []LabelRequirement {
	if len(requirements) == 0 {
		return nil
	}

	out := make([]LabelRequirement, len(requirements))
	for i, req := range requirements {
		out[i] = LabelRequirement{req: req.clone()}
	}

	return out
}

// matchLabelRequirements applies AND semantics without allocating per item.
func matchLabelRequirements(requirements []LabelRequirement, lookup func(string) (string, bool)) bool {
	for _, req := range requirements {
		if !req.req.match(lookup) {
			return false
		}
	}

	return true
}
