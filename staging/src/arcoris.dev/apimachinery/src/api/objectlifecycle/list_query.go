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

package objectlifecycle

import (
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resource"
)

// compileListQuery validates and canonicalizes lifecycle list filtering.
func compileListQuery(query objectquery.Query) (objectquery.Predicate, error) {
	predicate, err := objectquery.Compile(query)
	if err != nil {
		return objectquery.Predicate{}, errorFor(
			OperationList,
			ErrorReasonInvalidQuery,
			objectstore.Key{},
			ErrInvalidRequest,
			err,
		)
	}

	return predicate, nil
}

// validateListQueryForResourceAndScope rejects contradictory identity filters
// before storage is read.
func validateListQueryForResourceAndScope(
	resolved resolvedResource,
	scope objectstore.ListScope,
	predicate objectquery.Predicate,
) error {
	requirement := predicate.NamespaceRequirement()
	switch requirement.Kind {
	case objectquery.NamespaceUnconstrained,
		objectquery.NamespaceDisjunctive,
		objectquery.NamespaceResidual:
		return nil
	case objectquery.NamespaceContradictory:
		return invalidListQuery()
	}

	namespace := requirement.Namespace
	if scope.IsNamespace() && namespace != scope.Namespace() {
		return invalidListQuery()
	}

	switch resolved.definition.Scope() {
	case resource.ScopeGlobal:
		if namespace.IsZero() {
			return nil
		}
		return invalidListQuery()
	case resource.ScopeNamespaced:
		if namespace.IsZero() {
			return invalidListQuery()
		}
		return nil
	default:
		return errorFor(
			OperationList,
			ErrorReasonInvalidResourceContract,
			objectstore.Key{},
			ErrValidationFailed,
			ErrInvalidResourceContract,
		)
	}
}

// invalidListQuery maps query/scope/resource contradictions to lifecycle input
// errors without blaming objectstore or objectquery.
func invalidListQuery() error {
	return errorFor(
		OperationList,
		ErrorReasonInvalidQueryScope,
		objectstore.Key{},
		ErrInvalidRequest,
		nil,
	)
}
