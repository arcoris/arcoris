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
	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectvalidation"
	"arcoris.dev/apimachinery/api/value"
)

// validateObject delegates descriptor-aware object checks to objectvalidation.
func (e *Executor) validateObject(
	op Operation,
	key objectstore.Key,
	obj objectapply.ValueObject,
	resolved resolvedResource,
) error {
	plan := objectvalidation.Plan[value.Value, value.Value]{
		Resource:          resolved.definition,
		Resolver:          e.resolver,
		DesiredValidator:  e.desiredValidator,
		ObservedValidator: e.observedValidator,
	}

	if err := objectvalidation.Validate(obj, plan); err != nil {
		return errorFor(op, ReasonValidationFailed, key, ErrValidationFailed, err)
	}

	return nil
}
