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

package valuecompare

import "arcoris.dev/apimachinery/api/fieldpath"

// ValidateStructure checks result bucket well-formedness.
//
// It does not validate either compared value, descriptor semantics, apply
// policy, ownership state, or conflict behavior.
func (r Result) ValidateStructure() error {
	if err := validateResultBucket("added", r.added); err != nil {
		return err
	}
	if err := validateResultBucket("removed", r.removed); err != nil {
		return err
	}
	if err := validateResultBucket("modified", r.modified); err != nil {
		return err
	}

	if err := rejectOverlappingResultBuckets("added", r.added, "removed", r.removed); err != nil {
		return err
	}
	if err := rejectOverlappingResultBuckets("added", r.added, "modified", r.modified); err != nil {
		return err
	}
	if err := rejectOverlappingResultBuckets("removed", r.removed, "modified", r.modified); err != nil {
		return err
	}

	return nil
}

// validateResultBucket checks one comparison bucket for structural path validity.
func validateResultBucket(name string, set fieldpath.Set) error {
	var resultErr error
	set.ForEach(func(_ int, path fieldpath.Path) bool {
		if err := path.ValidateStructure(); err != nil {
			resultErr = wrapAt(
				path,
				ErrInvalidResult,
				ErrorReasonInvalidResultBucket,
				"result "+name+" bucket contains an invalid field path",
				err,
			)
			return false
		}

		return true
	})

	return resultErr
}

// rejectOverlappingResultBuckets rejects the same path appearing in two buckets.
func rejectOverlappingResultBuckets(
	leftName string,
	left fieldpath.Set,
	rightName string,
	right fieldpath.Set,
) error {
	var resultErr error
	left.ForEach(func(_ int, path fieldpath.Path) bool {
		if !right.Has(path) {
			return true
		}

		resultErr = errorfAt(
			path,
			ErrInvalidResult,
			ErrorReasonOverlappingResultPath,
			"result path appears in both %s and %s buckets",
			leftName,
			rightName,
		)
		return false
	})

	return resultErr
}
