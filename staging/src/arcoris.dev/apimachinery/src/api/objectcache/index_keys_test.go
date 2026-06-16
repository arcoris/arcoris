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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
)

func TestIndexKeysAreComparableAndDistinct(t *testing.T) {
	objects := map[objectNameKey]string{
		{namespace: metaidentity.Namespace("system"), name: metaidentity.Name("worker")}: "namespaced",
		{namespace: "", name: metaidentity.Name("worker")}:                               "global",
	}
	if got := objects[objectNameKey{namespace: "", name: "worker"}]; got != "global" {
		t.Fatalf("objectNameKey lookup = %q; want global", got)
	}

	labelsByValue := map[labelValueKey]string{
		{key: labels.Key("env"), value: labels.Value("prod")}: "prod",
		{key: labels.Key("env"), value: labels.Value("qa")}:   "qa",
	}
	if got := labelsByValue[labelValueKey{key: "env", value: "prod"}]; got != "prod" {
		t.Fatalf("labelValueKey lookup = %q; want prod", got)
	}

	annotationsByValue := map[annotationValueKey]string{
		{key: annotations.Key("team"), value: annotations.Value("core")}: "core",
		{key: annotations.Key("team"), value: annotations.Value("ops")}:  "ops",
	}
	if got := annotationsByValue[annotationValueKey{key: "team", value: "ops"}]; got != "ops" {
		t.Fatalf("annotationValueKey lookup = %q; want ops", got)
	}
}
