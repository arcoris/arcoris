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

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestCompareNullBothNullIsEmpty(t *testing.T) {
	got, err := newComparer(Options{}).compareNull(rootField("name"), value.NullValue(), value.NullValue())
	requireNoError(t, err)

	requireResult(t, got, nil, nil, nil)
}

func TestCompareNullScalarChangeIsModified(t *testing.T) {
	path := rootField("name")

	got, err := newComparer(Options{}).compareNull(path, value.NullValue(), value.StringValue("x"))
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareNullDescriptorRejectsNonNull(t *testing.T) {
	_, err := newComparer(Options{}).compareNullDescriptor(
		rootField("name"),
		value.StringValue("x"),
		value.NullValue(),
		types.Null().Type(),
	)

	requireErrorIs(t, err, ErrKindMismatch)
}
