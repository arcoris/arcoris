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

package apidocument_test

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestPathsTree(t *testing.T) {
	paths := apidocument.Paths()

	tests := []struct {
		name string
		got  apidocument.Path
		want string
	}{
		{name: "typeMeta apiVersion", got: paths.TypeMeta().APIVersion(), want: "typeMeta.apiVersion"},
		{name: "typeMeta kind", got: paths.TypeMeta().Kind(), want: "typeMeta.kind"},
		{name: "object apiVersion", got: paths.Object().APIVersion(), want: "object.apiVersion"},
		{name: "object kind", got: paths.Object().Kind(), want: "object.kind"},
		{name: "object metadata", got: paths.Object().MetadataPath(), want: "object.metadata"},
		{name: "object metadata labels", got: paths.Object().Metadata().Labels(), want: "object.metadata.labels"},
		{name: "object metadata annotations", got: paths.Object().Metadata().Annotations(), want: "object.metadata.annotations"},
		{name: "object desired", got: paths.Object().Desired(), want: "object.desired"},
		{name: "object observed", got: paths.Object().Observed(), want: "object.observed"},
		{name: "objectMeta name", got: paths.ObjectMeta().Name(), want: "objectMeta.name"},
		{name: "objectMeta labels", got: paths.ObjectMeta().Labels(), want: "objectMeta.labels"},
		{name: "ownership desired", got: paths.Ownership().Desired(), want: "ownership.desired"},
		{name: "ownership desired entries", got: paths.Ownership().DesiredSurface().Entries(), want: "ownership.desired.entries"},
		{name: "ownership desired entry owner", got: paths.Ownership().DesiredSurface().Entry().Owner(), want: "ownership.desired.entries.owner"},
		{name: "ownership observed", got: paths.Ownership().Observed(), want: "ownership.observed"},
		{name: "ownership metadata", got: paths.Ownership().MetadataPath(), want: "ownership.metadata"},
		{name: "ownership metadata labels", got: paths.Ownership().Metadata().Labels(), want: "ownership.metadata.labels"},
		{name: "ownership metadata labels entries", got: paths.Ownership().Metadata().LabelsSurface().Entries(), want: "ownership.metadata.labels.entries"},
		{name: "ownership metadata annotations", got: paths.Ownership().Metadata().Annotations(), want: "ownership.metadata.annotations"},
		{name: "pageMeta resourceVersion", got: paths.PageMeta().ResourceVersion(), want: "pageMeta.resourceVersion"},
		{name: "pageMeta continue", got: paths.PageMeta().Continue(), want: "pageMeta.continue"},
		{name: "pageMeta remainingItemCount", got: paths.PageMeta().RemainingItemCount(), want: "pageMeta.remainingItemCount"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.String() != tt.want {
				t.Fatalf("path = %q; want %q", tt.got, tt.want)
			}
		})
	}
}
