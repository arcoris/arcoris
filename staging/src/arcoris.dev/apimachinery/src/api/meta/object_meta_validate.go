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

package meta

// Validate checks metadata-level invariants only.
//
// Scope-specific rules such as requiring or forbidding Namespace belong to
// resource-aware validation outside api/meta.
func (m ObjectMeta) Validate() error {
	if !m.Name.IsZero() {
		if err := m.Name.Validate(); err != nil {
			return nested("objectMeta.name", ErrInvalidObjectMeta, err)
		}
	}

	if !m.GenerateName.IsZero() {
		if err := m.GenerateName.Validate(); err != nil {
			return nested("objectMeta.generateName", ErrInvalidObjectMeta, err)
		}
	}

	if !m.Namespace.IsZero() {
		if err := m.Namespace.Validate(); err != nil {
			return nested("objectMeta.namespace", ErrInvalidObjectMeta, err)
		}
	}

	if !m.UID.IsZero() {
		if err := m.UID.Validate(); err != nil {
			return nested("objectMeta.uid", ErrInvalidObjectMeta, err)
		}
	}

	if !m.ResourceVersion.IsZero() {
		if err := m.ResourceVersion.Validate(); err != nil {
			return nested("objectMeta.resourceVersion", ErrInvalidObjectMeta, err)
		}
	}

	if err := m.Generation.Validate(); err != nil {
		return nested("objectMeta.generation", ErrInvalidObjectMeta, err)
	}

	if !m.CreatedAt.IsZero() {
		if err := m.CreatedAt.Validate(); err != nil {
			return nested("objectMeta.createdAt", ErrInvalidObjectMeta, err)
		}
	}

	if m.Deletion != nil {
		if err := m.Deletion.Validate(); err != nil {
			return nested("objectMeta.deletion", ErrInvalidObjectMeta, err)
		}
	}

	if err := m.Labels.Validate(); err != nil {
		return nested("objectMeta.labels", ErrInvalidObjectMeta, err)
	}

	if err := m.Annotations.Validate(); err != nil {
		return nested("objectMeta.annotations", ErrInvalidObjectMeta, err)
	}

	if err := m.OwnerReferences.Validate(); err != nil {
		return nested("objectMeta.ownerReferences", ErrInvalidObjectMeta, err)
	}

	if err := m.Finalizers.Validate(); err != nil {
		return nested("objectMeta.finalizers", ErrInvalidObjectMeta, err)
	}

	return nil
}
