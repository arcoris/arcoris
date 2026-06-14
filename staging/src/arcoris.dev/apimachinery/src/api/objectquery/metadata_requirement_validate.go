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

package objectquery

import (
	"errors"
	"strconv"
)

// validate checks operator validity, operator/value arity, key syntax, and
// value syntax for one metadata requirement.
func (r metadataRequirement) validate(path string, validateKey validatorFunc, validateValue validatorFunc) error {
	if !r.op.IsValid() {
		return errorf(
			path+".operator",
			errors.Join(ErrInvalidRequirement, ErrInvalidOperator),
			ErrorReasonInvalidOperator,
			"operator %q is invalid",
			r.op,
		)
	}
	if err := validateKey(r.key); err != nil {
		return wrapf(
			path+".key",
			ErrInvalidRequirement,
			ErrorReasonInvalidKey,
			err,
			"metadata key is invalid",
		)
	}
	if err := validateMetadataValueCount(path+".values", r.op, len(r.values)); err != nil {
		return err
	}
	for i, value := range r.values {
		if err := validateValue(value); err != nil {
			return wrapf(
				path+".values["+strconv.Itoa(i)+"]",
				ErrInvalidRequirement,
				ErrorReasonInvalidValue,
				err,
				"metadata value is invalid",
			)
		}
	}

	return nil
}

// validateMetadataValueCount enforces the finite operator arity contract before
// values are sorted and deduplicated.
func validateMetadataValueCount(path string, op Operator, count int) error {
	switch op {
	case OperatorExists, OperatorDoesNotExist:
		if count != 0 {
			return errorf(path, ErrInvalidRequirement, ErrorReasonInvalidValueCount, "%s requires no values", op)
		}
	case OperatorEquals, OperatorNotEquals:
		if count != 1 {
			return errorf(path, ErrInvalidRequirement, ErrorReasonInvalidValueCount, "%s requires exactly one value", op)
		}
	case OperatorIn, OperatorNotIn:
		if count == 0 {
			return errorf(path, ErrInvalidRequirement, ErrorReasonInvalidValueCount, "%s requires at least one value", op)
		}
	default:
		return errorf(path, errors.Join(ErrInvalidRequirement, ErrInvalidOperator), ErrorReasonInvalidOperator, "operator %q is invalid", op)
	}

	return nil
}
