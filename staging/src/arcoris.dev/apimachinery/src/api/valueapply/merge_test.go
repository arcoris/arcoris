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

	"arcoris.dev/apimachinery/api/types"
)

func TestMerge(t *testing.T) {
	req := specRequest(owner("user"))
	result := Result{MergeFields: fields(imagePath())}

	got, err := newApplier(Options{}).merge(req, result)
	requireNoError(t, err)

	requireStringMember(t, got, "image", "api:v2")
}

func TestMergeOptions(t *testing.T) {
	got := newApplier(Options{MaxDepth: 13}).mergeOptions()

	if got.MaxDepth != 13 {
		t.Fatalf("MaxDepth = %d; want 13", got.MaxDepth)
	}
}

func TestApplyValueMergeErrorWrapped(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       str("old"),
		Applied:    str("new"),
		Descriptor: types.String().Type(),
	}
	result := Result{MergeFields: fields(root().Field("nested"))}

	_, err := newApplier(Options{}).merge(req, result)

	requireErrorIs(t, err, ErrMergeFailed)
}
