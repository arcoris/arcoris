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

// ValidateDefinitionLocal checks the local structural integrity of def.
//
// Local validation checks identity values, scope, version-set invariants, and
// descriptor-local Desired/Observed surfaces. It validates reference syntax but
// does not require referenced descriptor definitions to resolve.
func ValidateDefinitionLocal(def Definition) error {
	if err := validateDefinitionIdentity(def); err != nil {
		return err
	}
	return validateDefinitionVersions(def.versions, nil)
}

// ValidateDefinitionResolved checks the structural integrity of def.
//
// Resolved validation checks identity values, scope, version-set invariants, and
// Desired/Observed structural surfaces, and it requires descriptor references to
// resolve through resolver. It does not validate concrete object values, derive
// metadata, export schemas, perform conversion/defaulting, define storage
// behavior, or register the definition in a global catalog.
func ValidateDefinitionResolved(def Definition, resolver types.Resolver) error {
	if err := validateDefinitionIdentity(def); err != nil {
		return err
	}
	return validateDefinitionVersions(def.versions, resolver)
}
