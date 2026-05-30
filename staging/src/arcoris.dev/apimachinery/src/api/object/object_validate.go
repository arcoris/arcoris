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

package object

// ValidateMeta validates only TypeMeta and ObjectMeta.
//
// Desired and observed payloads are intentionally ignored. Payload validation,
// scope checks, and surface descriptor checks require resource-aware context
// outside api/object.
func (o Object[D, O]) ValidateMeta() error {
	if err := o.TypeMeta.Validate(); err != nil {
		return nested("object.typeMeta", ErrInvalidObject, err)
	}

	if err := o.ObjectMeta.Validate(); err != nil {
		return nested("object.metadata", ErrInvalidObject, err)
	}

	return nil
}
