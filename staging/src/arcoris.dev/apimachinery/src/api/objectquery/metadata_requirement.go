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
	"slices"
	"strconv"
	"strings"
)

// metadataRequirement is the canonical shared representation for labels and annotations.
type metadataRequirement struct {
	key    string
	op     Operator
	values []string
}

// validatorFunc adapts label and annotation lexical validators.
type validatorFunc func(string) error

// metadataRequirementFrom constructs and canonicalizes one metadata requirement.
func metadataRequirementFrom(
	path string,
	key string,
	op Operator,
	values []string,
	validateKey validatorFunc,
	validateValue validatorFunc,
) (metadataRequirement, error) {
	req := metadataRequirement{
		key:    key,
		op:     op,
		values: append([]string(nil), values...),
	}
	if err := req.validate(path, validateKey, validateValue); err != nil {
		return metadataRequirement{}, err
	}
	req.values = canonicalValues(req.values)

	return req, nil
}

// validate checks requirement shape, operator arity, and metadata lexical rules.
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
	if err := validateValueCount(path+".values", r.op, len(r.values)); err != nil {
		return err
	}
	for i, value := range r.values {
		if err := validateValue(value); err != nil {
			return wrapf(
				path+".values["+itoa(i)+"]",
				ErrInvalidRequirement,
				ErrorReasonInvalidValue,
				err,
				"metadata value is invalid",
			)
		}
	}

	return nil
}

// validateValueCount checks operator/value arity before value canonicalization.
func validateValueCount(path string, op Operator, count int) error {
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

// canonicalValues sorts and deduplicates membership values.
func canonicalValues(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	out := append([]string(nil), values...)
	slices.Sort(out)
	return slices.Compact(out)
}

// clone returns a detached requirement.
func (r metadataRequirement) clone() metadataRequirement {
	r.values = append([]string(nil), r.values...)
	return r
}

// match evaluates r against one metadata key/value lookup function.
func (r metadataRequirement) match(lookup func(string) (string, bool)) bool {
	actual, ok := lookup(r.key)
	switch r.op {
	case OperatorExists:
		return ok
	case OperatorDoesNotExist:
		return !ok
	case OperatorEquals:
		return ok && actual == r.values[0]
	case OperatorNotEquals:
		return !ok || actual != r.values[0]
	case OperatorIn:
		return ok && slices.Contains(r.values, actual)
	case OperatorNotIn:
		return !ok || !slices.Contains(r.values, actual)
	default:
		return false
	}
}

// compareMetadataRequirement orders requirements by key, operator, then values.
func compareMetadataRequirement(left, right metadataRequirement) int {
	if cmp := strings.Compare(left.key, right.key); cmp != 0 {
		return cmp
	}
	if left.op < right.op {
		return -1
	}
	if left.op > right.op {
		return 1
	}

	return compareStringSlices(left.values, right.values)
}

// sameMetadataRequirement reports exact canonical requirement equality.
func sameMetadataRequirement(left, right metadataRequirement) bool {
	return left.key == right.key &&
		left.op == right.op &&
		slices.Equal(left.values, right.values)
}

// compareStringSlices orders string slices lexicographically by content.
func compareStringSlices(left, right []string) int {
	for i := 0; i < len(left) && i < len(right); i++ {
		if cmp := strings.Compare(left[i], right[i]); cmp != 0 {
			return cmp
		}
	}
	switch {
	case len(left) < len(right):
		return -1
	case len(left) > len(right):
		return 1
	default:
		return 0
	}
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
