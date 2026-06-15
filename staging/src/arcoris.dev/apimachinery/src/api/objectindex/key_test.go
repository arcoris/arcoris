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
	"arcoris.dev/apimachinery/api/meta/labels"
)

func TestIndexKeysAreComparableMapKeys(t *testing.T) {
	objectNames := map[objectNameKey]int{
		{namespace: "system", name: "worker-1"}: 1,
	}
	labelValues := map[labelValueKey]int{
		{key: labels.Key("env"), value: labels.Value("prod")}: 2,
	}
	annotationValues := map[annotationValueKey]int{
		{key: annotations.Key("team"), value: annotations.Value("core")}: 3,
	}

	if got := objectNames[objectNameKey{namespace: "system", name: "worker-1"}]; got != 1 {
		t.Fatalf("objectNameKey lookup = %d; want 1", got)
	}
	if got := labelValues[labelValueKey{key: labels.Key("env"), value: labels.Value("prod")}]; got != 2 {
		t.Fatalf("labelValueKey lookup = %d; want 2", got)
	}
	if got := annotationValues[annotationValueKey{key: annotations.Key("team"), value: annotations.Value("core")}]; got != 3 {
		t.Fatalf("annotationValueKey lookup = %d; want 3", got)
	}
}
