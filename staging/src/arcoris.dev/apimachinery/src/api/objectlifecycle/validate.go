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
	"context"
	"reflect"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
)

// validateConfig checks constructor dependencies that every operation needs.
func validateConfig(cfg config) error {
	switch {
	case isNilInterface(cfg.store):
		return errorFor(0, ErrorReasonInvalidExecutor, objectstore.Key{}, ErrInvalidExecutor, ErrNilStore)
	case isNilInterface(cfg.resources):
		return errorFor(0, ErrorReasonInvalidExecutor, objectstore.Key{}, ErrInvalidExecutor, ErrNilResourceResolver)
	case isNilInterface(cfg.desiredValidator):
		return errorFor(0, ErrorReasonInvalidExecutor, objectstore.Key{}, ErrInvalidExecutor, ErrNilDesiredValidator)
	default:
		return nil
	}
}

// requireExecutor rejects nil or partially initialized executors.
func (e *Executor) requireExecutor(op Operation) error {
	if e == nil ||
		isNilInterface(e.store) ||
		isNilInterface(e.resources) ||
		isNilInterface(e.desiredValidator) {
		return errorFor(op, ErrorReasonInvalidExecutor, objectstore.Key{}, ErrInvalidExecutor, ErrInvalidExecutor)
	}

	return nil
}

// checkContext rejects nil contexts before lower layers dereference them.
func checkContext(op Operation, ctx context.Context) error {
	if ctx == nil {
		return errorFor(op, ErrorReasonInvalidContext, objectstore.Key{}, ErrInvalidRequest, ErrNilContext)
	}

	return nil
}

// validateOwner checks the field manager identity before apply or ownership init.
func validateOwner(op Operation, owner fieldownership.Owner) error {
	if err := owner.ValidateLexical(); err != nil {
		return errorFor(op, ErrorReasonInvalidOwner, objectstore.Key{}, ErrInvalidRequest, err)
	}

	return nil
}

// validateCreateRequest checks lifecycle-local Create request fields.
func (e *Executor) validateCreateRequest(req CreateRequest) error {
	return validateOwner(OperationCreate, req.Owner)
}

// validateApplyRequest checks lifecycle-local Apply request fields.
func (e *Executor) validateApplyRequest(req ApplyRequest) error {
	if err := validateOwner(OperationApply, req.Owner); err != nil {
		return err
	}
	if req.Object.Observed != nil {
		return errorFor(
			OperationApply,
			ErrorReasonUnsupportedObservedApply,
			objectstore.Key{},
			ErrInvalidRequest,
			objectapply.ErrUnsupportedObservedApply,
		)
	}

	return nil
}

// validateGetRequest exists to make the read pipeline stage explicit.
func (e *Executor) validateGetRequest(GetRequest) error {
	return nil
}

// validateDeleteRequest checks lifecycle-local Delete request fields.
func (e *Executor) validateDeleteRequest(req DeleteRequest) error {
	if req.Expected.IsZero() {
		return errorFor(
			OperationDelete,
			ErrorReasonInvalidExpectedRevision,
			objectstore.Key{},
			ErrInvalidRequest,
			objectstore.ErrInvalidRevision,
		)
	}

	return nil
}

// validateUpdateObservedRequest checks lifecycle-local UpdateObserved request fields.
func (e *Executor) validateUpdateObservedRequest(req UpdateObservedRequest) error {
	if err := validateOwner(OperationUpdateObserved, req.Owner); err != nil {
		return err
	}
	if req.Expected.IsZero() {
		return errorFor(
			OperationUpdateObserved,
			ErrorReasonInvalidExpectedRevision,
			objectstore.Key{},
			ErrInvalidRequest,
			objectstore.ErrInvalidRevision,
		)
	}
	if isNilInterface(e.observedValidator) {
		return errorFor(
			OperationUpdateObserved,
			ErrorReasonInvalidExecutor,
			objectstore.Key{},
			ErrInvalidExecutor,
			ErrNilObservedValidator,
		)
	}

	return nil
}

// validatePatchMetadataRequest checks lifecycle-local PatchMetadata request fields.
func (e *Executor) validatePatchMetadataRequest(req PatchMetadataRequest) error {
	if err := validateOwner(OperationPatchMetadata, req.Owner); err != nil {
		return err
	}
	if req.Expected.IsZero() {
		return errorFor(
			OperationPatchMetadata,
			ErrorReasonInvalidExpectedRevision,
			objectstore.Key{},
			ErrInvalidRequest,
			objectstore.ErrInvalidRevision,
		)
	}
	if len(req.Labels) == 0 && len(req.Annotations) == 0 {
		return errorFor(
			OperationPatchMetadata,
			ErrorReasonEmptyMetadataPatch,
			objectstore.Key{},
			ErrInvalidRequest,
			nil,
		)
	}

	return nil
}

// validateExpectedRevision maps store revision validation into lifecycle detail.
func validateExpectedRevision(op Operation, key objectstore.Key, expected objectstore.Revision) error {
	if err := objectstore.ValidateExpectedRevision(key, expected); err != nil {
		return errorFor(op, ErrorReasonInvalidExpectedRevision, key, ErrInvalidRequest, err)
	}

	return nil
}

// isNilInterface detects nil and typed-nil interface values.
func isNilInterface(value any) bool {
	if value == nil {
		return true
	}

	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return reflectValue.IsNil()
	default:
		return false
	}
}
