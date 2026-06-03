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

func TestAddSubtreeUsesValueFieldSet(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Type()
	val := value.MustObjectValue(value.ObjectMember("image", value.StringValue("v1")))

	got, err := newComparer(Options{}).addSubtree(path, val, descriptor, EmptyResult())
	requireNoError(t, err)

	requireResult(t, got, paths(path.Field("image")), nil, nil)
}

func TestRemoveSubtreeUsesValueFieldSet(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Type()
	val := value.MustObjectValue(value.ObjectMember("image", value.StringValue("v1")))

	got, err := newComparer(Options{}).removeSubtree(path, val, descriptor, EmptyResult())
	requireNoError(t, err)

	requireResult(t, got, nil, paths(path.Field("image")), nil)
}
