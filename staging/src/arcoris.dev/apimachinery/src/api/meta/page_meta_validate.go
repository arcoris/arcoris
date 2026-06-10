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

// ValidateLexical checks page metadata scalar form only.
//
// It does not interpret pagination policy, storage cursor internals, list
// consistency guarantees, watch state, or remaining-count accuracy.
func (m PageMeta) ValidateLexical() error {
	if !m.ResourceVersion.IsZero() {
		if err := m.ResourceVersion.ValidateLexical(); err != nil {
			return nested("pageMeta.resourceVersion", ErrInvalidPageMeta, err)
		}
	}

	if !m.ContinueToken.IsZero() {
		if err := m.ContinueToken.ValidateLexical(); err != nil {
			return nested("pageMeta.continue", ErrInvalidPageMeta, err)
		}
	}

	return nil
}
