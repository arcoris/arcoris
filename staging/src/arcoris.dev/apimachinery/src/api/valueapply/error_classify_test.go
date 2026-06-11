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

package valueapply

import (
	"testing"

	"arcoris.dev/apimachinery/api/valuecompare"
	"arcoris.dev/apimachinery/api/valuefieldset"
	"arcoris.dev/apimachinery/api/valuemerge"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestWrapValidationErrorClassifiesDescriptorFailures(t *testing.T) {
	err := wrapValidationError(root(), "live validation failed", valuevalidation.ErrInvalidDescriptor)

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorIs(t, err, valuevalidation.ErrInvalidDescriptor)
}

func TestWrapFieldSetErrorClassifiesPathFailures(t *testing.T) {
	err := wrapFieldSetError(root(), valuefieldset.ErrInvalidPath)

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, valuefieldset.ErrInvalidPath)
}

func TestWrapCompareErrorClassifiesReferenceFailures(t *testing.T) {
	err := wrapCompareError(root(), valuecompare.ErrReferenceCycle)

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorIs(t, err, valuecompare.ErrReferenceCycle)
}

func TestWrapMergeErrorClassifiesUnsupportedMerge(t *testing.T) {
	err := wrapMergeError(root(), valuemerge.ErrUnsupportedMerge)

	requireErrorIs(t, err, ErrUnsupportedMerge)
	requireErrorIs(t, err, valuemerge.ErrUnsupportedMerge)
}
