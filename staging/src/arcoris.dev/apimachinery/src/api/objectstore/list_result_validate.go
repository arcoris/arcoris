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

package objectstore

import "errors"

// ValidateListResult checks that result is a structurally valid response for request.
//
// A zero result Revision is valid for an empty collection boundary: it means no
// committed mutation has been observed yet. Committed item revisions remain
// non-zero and must not be newer than the result watermark.
func ValidateListResult(request ListRequest, result ListResult) error {
	if err := ValidateListRequest(request); err != nil {
		return errorFor(
			ErrorReasonInvalidListResult,
			Key{},
			0,
			0,
			errors.Join(ErrInvalidListResult, err),
		)
	}

	for _, item := range result.Items {
		if err := ValidateListItem(item); err != nil {
			return errorFor(
				ErrorReasonInvalidListItem,
				item.Key,
				0,
				0,
				errors.Join(ErrInvalidListResult, err),
			)
		}
		if !KeyMatchesListRequest(item.Key, request) {
			return errorFor(
				ErrorReasonInvalidListItem,
				item.Key,
				0,
				0,
				errors.Join(ErrInvalidListResult, ErrInvalidListRequest),
			)
		}
		if result.Revision.Before(item.State.Revision) {
			return errorFor(
				ErrorReasonInvalidListResult,
				item.Key,
				result.Revision,
				item.State.Revision,
				errors.Join(ErrInvalidListResult, ErrInvalidRevision),
			)
		}
	}

	return nil
}
