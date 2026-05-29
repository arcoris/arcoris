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

// Descriptor paths used in structured diagnostics.
//
// Keeping paths in one place prevents validators and encoding code from
// drifting into subtly different names for the same descriptor location.
const (
	pathDefinitionGroup    = "definition.group"
	pathDefinitionKind     = "definition.kind"
	pathDefinitionResource = "definition.resource"
	pathDefinitionScope    = "definition.scope"
	pathDefinitionVersions = "definition.versions"

	pathScope = "scope"
)

// Scope grammar tokens.
//
// Scope has a deliberately tiny grammar because the resource package only
// describes API resource families. It does not define routing, namespace value
// objects, tenancy, or storage partitioning behavior.
const (
	scopeTextGlobal     = "global"
	scopeTextNamespaced = "namespaced"
	scopeTextInvalid    = "invalid"
	scopeTextUnknown    = "unknown"
)

// Reusable diagnostic details for scope validation and parsing.
//
// The text is centralized so parse, validate, and encoding paths report the
// same human-facing contract.
const (
	detailScopeRequired  = "scope must be non-empty"
	detailScopeSupported = "scope must be global or namespaced"
)

// Reusable definition-level diagnostic details.
const (
	detailDefinitionNeedsVersion     = "resource definition must contain at least one version"
	detailDefinitionNeedsExposed     = "resource definition must expose at least one version"
	detailDefinitionNeedsCanonical   = "resource definition must declare exactly one canonical version"
	detailDefinitionCanonicalExposed = "canonical version must also be exposed"
	detailVersionDesiredRequired     = "desired descriptor must be present"
	detailDesiredObjectLikeTemplate  = "desired"
	detailObservedObjectLikeTemplate = "observed"
	detailJSONMustBeString           = "JSON value must be a string"
	detailJSONMustBeNonNullString    = "JSON value must be a non-null string"
	detailDecodeTargetMustBeNonNil   = "decode target must be non-nil"
)
