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

	"arcoris.dev/apimachinery/api/meta/labels"
)

func TestIndexLabelsRecordsExistenceAndValuePositions(t *testing.T) {
	idx := Build(testItems())

	requirePositions(t, idx.byLabelKey[labels.Key("env")], 0, 1, 2, 3)
	requirePositions(t, idx.byLabelKey[labels.Key("tier")], 0, 1, 2)
	requirePositions(t, idx.byLabelValue[labelValueKey{
		key:   labels.Key("env"),
		value: labels.Value("prod"),
	}], 0, 2, 3)
	requirePositions(t, idx.byLabelValue[labelValueKey{
		key:   labels.Key("tier"),
		value: labels.Value("backend"),
	}], 0, 1)
}

func TestIndexLabelsOmitsItemsWithoutLabels(t *testing.T) {
	idx := Build(testItems())

	if got := idx.byLabelKey[labels.Key("missing")]; got != nil {
		t.Fatalf("missing label positions = %v; want nil", got)
	}
	if got := idx.byLabelValue[labelValueKey{key: labels.Key("env"), value: labels.Value("missing")}]; got != nil {
		t.Fatalf("missing label value positions = %v; want nil", got)
	}
}
