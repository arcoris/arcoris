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

package objectownership

import "arcoris.dev/apimachinery/api/fieldownership"

// NewMetadataState constructs metadata ownership state.
//
// Labels and annotations are independent ownership namespaces. A key path in
// labels never overlaps the same key path in annotations.
func NewMetadataState(labels fieldownership.State, annotations fieldownership.State) MetadataState {
	return MetadataState{labels: labels, annotations: annotations}
}

// IsEmpty reports whether no modeled metadata surface has ownership state.
func (m MetadataState) IsEmpty() bool {
	return m.labels.IsEmpty() &&
		m.annotations.IsEmpty()
}

// Labels returns metadata.labels ownership.
//
// Paths are relative to the labels map root. For example,
// $["scheduler.arcoris.dev/mode"] owns one label key and does not include a
// synthetic $.metadata.labels prefix.
func (m MetadataState) Labels() fieldownership.State {
	return m.labels
}

// Annotations returns metadata.annotations ownership.
//
// Paths are relative to the annotations map root, independent from labels.
func (m MetadataState) Annotations() fieldownership.State {
	return m.annotations
}

// WithLabels returns a copy of m with replacement labels ownership.
//
// Annotations ownership is preserved.
func (m MetadataState) WithLabels(labels fieldownership.State) MetadataState {
	m.labels = labels

	return m
}

// WithAnnotations returns a copy of m with replacement annotations ownership.
//
// Labels ownership is preserved.
func (m MetadataState) WithAnnotations(annotations fieldownership.State) MetadataState {
	m.annotations = annotations

	return m
}
