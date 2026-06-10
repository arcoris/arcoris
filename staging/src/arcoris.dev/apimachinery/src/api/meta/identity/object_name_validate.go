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

package identity

// ValidateLexical checks that an object name has a name and optional valid namespace.
func (n ObjectName) ValidateLexical() error {
	if n.Name.IsZero() {
		return invalid(
			"objectName.name",
			n.String(),
			ErrInvalidObjectName,
			ErrorReasonEmptyValue,
			"name is required",
		)
	}

	if err := n.Namespace.ValidateLexical(); err != nil {
		return nested("objectName.namespace", n.String(), ErrInvalidObjectName, err)
	}

	if err := n.Name.ValidateLexical(); err != nil {
		return nested("objectName.name", n.String(), ErrInvalidObjectName, err)
	}

	return nil
}
