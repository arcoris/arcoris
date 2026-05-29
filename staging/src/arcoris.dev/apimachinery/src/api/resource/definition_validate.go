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

package resource

import "arcoris.dev/apimachinery/api/types"

// ValidateDefinition checks the structural integrity of def.
//
// Validation is descriptor validation only. It checks identity values, scope,
// version-set invariants, and Desired/Observed structural surfaces. It does not
// validate concrete object values, derive metadata, export schemas, perform
// conversion/defaulting, define storage behavior, or register the definition in
// a global catalog.
func ValidateDefinition(def Definition, resolver types.Resolver) error {
	if err := validateDefinitionIdentity(def); err != nil {
		return err
	}
	return validateDefinitionVersions(def.versions, resolver)
}
