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

package objectindex

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta/annotations"
)

func TestIndexAnnotationsRecordsExistenceAndValuePositions(t *testing.T) {
	idx := Build(testItems())

	requirePositions(t, idx.byAnnotationKey[annotations.Key("team")], 0, 1, 2, 3)
	requirePositions(t, idx.byAnnotationKey[annotations.Key("zone")], 0, 1, 2)
	requirePositions(t, idx.byAnnotationValue[annotationValueKey{
		key:   annotations.Key("team"),
		value: annotations.Value("core"),
	}], 0, 2)
	requirePositions(t, idx.byAnnotationValue[annotationValueKey{
		key:   annotations.Key("zone"),
		value: annotations.Value("west"),
	}], 1, 2)
}

func TestIndexAnnotationsOmitsItemsWithoutAnnotations(t *testing.T) {
	idx := Build(testItems())

	if got := idx.byAnnotationKey[annotations.Key("missing")]; got != nil {
		t.Fatalf("missing annotation positions = %v; want nil", got)
	}
	if got := idx.byAnnotationValue[annotationValueKey{key: annotations.Key("team"), value: annotations.Value("missing")}]; got != nil {
		t.Fatalf("missing annotation value positions = %v; want nil", got)
	}
}
