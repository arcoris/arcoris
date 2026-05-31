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

package objectvalidation

import "fmt"

const (
	// pathPlanResource points at the resolved resource definition dependency.
	pathPlanResource = "plan.resource"

	// pathPlanDesiredValidator points at the required desired surface validator.
	pathPlanDesiredValidator = "plan.desiredValidator"

	// pathPlanObservedValidator points at the conditional observed surface validator.
	pathPlanObservedValidator = "plan.observedValidator"

	// pathObject points at the whole object envelope when a narrower path is unavailable.
	pathObject = "object"

	// pathObjectTypeMeta points at apiVersion/kind metadata.
	pathObjectTypeMeta = "object.typeMeta"

	// pathObjectNamespace points at the metadata namespace used for scope checks.
	pathObjectNamespace = "object.metadata.namespace"

	// pathObjectDesired points at the desired payload surface.
	pathObjectDesired = "object.desired"

	// pathObjectObserved points at the observed payload surface.
	pathObjectObserved = "object.observed"

	// pathResourceScope points at the resource definition scope.
	pathResourceScope = "resource.scope"

	// pathResourceVersions points at the resource version set.
	pathResourceVersions = "resource.versions"
)

// errorf builds a structured validation error with formatted detail text.
func errorf(path string, err error, reason ErrorReason, format string, args ...any) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: fmt.Sprintf(format, args...),
	}
}

// nested wraps a lower-layer validation failure without hiding its identity.
func nested(path string, err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
		Cause:  cause,
	}
}

// missingValidator reports a required typed surface validator dependency.
func missingValidator(path string, detail string) error {
	return &Error{
		Path:   path,
		Err:    ErrMissingValidator,
		Reason: ErrorReasonMissingValidator,
		Detail: detail,
	}
}
