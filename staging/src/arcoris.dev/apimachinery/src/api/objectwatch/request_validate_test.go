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

package objectwatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestRequestValidateRejectsInvalidListRequest(t *testing.T) {
	request := Request{
		List:  objectstore.ListRequest{Scope: objectstore.AllNamespaces()},
		Start: AtCurrent(),
	}

	err := request.Validate()

	requireErrorIs(t, err, ErrInvalidRequest)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
	requireWatchError(t, err, ErrorReasonInvalidRequest, "watch.request.list")
}

func TestRequestValidateRejectsInvalidStart(t *testing.T) {
	request := Request{List: watchListRequest(), Start: Start{}}

	err := request.Validate()

	requireErrorIs(t, err, ErrInvalidRequest)
	requireErrorIs(t, err, ErrInvalidStart)
	requireWatchError(t, err, ErrorReasonInvalidStart, "watch.request.start")
}

func TestRequestValidateAllowBookmarksIsNotStructural(t *testing.T) {
	for _, allowBookmarks := range []bool{false, true} {
		request := Request{
			List:           watchListRequest(),
			Start:          AtCurrent(),
			AllowBookmarks: allowBookmarks,
		}

		requireNoError(t, request.Validate())
	}
}
