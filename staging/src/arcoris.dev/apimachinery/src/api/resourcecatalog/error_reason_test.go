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

package resourcecatalog

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	requireEqual(t, string(ErrorReasonInvalidDefinition), "invalid_definition")
	requireEqual(t, string(ErrorReasonDuplicateResource), "duplicate_resource")
	requireEqual(t, string(ErrorReasonDuplicateKind), "duplicate_kind")
	requireEqual(t, string(ErrorReasonDuplicateVersionResource), "duplicate_version_resource")
	requireEqual(t, string(ErrorReasonDuplicateVersionKind), "duplicate_version_kind")
	requireEqual(t, string(ErrorReasonDefinitionExists), "definition_exists")
	requireEqual(t, string(ErrorReasonNilCatalog), "nil_catalog")
}
