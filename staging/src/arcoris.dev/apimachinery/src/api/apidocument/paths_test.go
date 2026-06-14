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

func TestPathsTreeHasNoDuplicatePaths(t *testing.T) {
	paths := exportedDocumentPaths()
	seen := make(map[apidocument.Path]string, len(paths))

	for name, path := range paths {
		if path.IsZero() {
			t.Fatalf("%s path is zero", name)
		}
		if previous, ok := seen[path]; ok {
			t.Fatalf("%s and %s both expose path %q", previous, name, path)
		}
		seen[path] = name
	}
}

func exportedDocumentPaths() map[string]apidocument.Path {
	paths := apidocument.Paths()

	return map[string]apidocument.Path{
		"typeMeta.apiVersion":                           paths.TypeMeta().APIVersion(),
		"typeMeta.kind":                                 paths.TypeMeta().Kind(),
		"object.apiVersion":                             paths.Object().APIVersion(),
		"object.kind":                                   paths.Object().Kind(),
		"object.metadata":                               paths.Object().MetadataPath(),
		"object.metadata.name":                          paths.Object().Metadata().Name(),
		"object.metadata.generateName":                  paths.Object().Metadata().GenerateName(),
		"object.metadata.namespace":                     paths.Object().Metadata().Namespace(),
		"object.metadata.uid":                           paths.Object().Metadata().UID(),
		"object.metadata.resourceVersion":               paths.Object().Metadata().ResourceVersion(),
		"object.metadata.generation":                    paths.Object().Metadata().Generation(),
		"object.metadata.createdAt":                     paths.Object().Metadata().CreatedAt(),
		"object.metadata.deletion":                      paths.Object().Metadata().Deletion(),
		"object.metadata.labels":                        paths.Object().Metadata().Labels(),
		"object.metadata.annotations":                   paths.Object().Metadata().Annotations(),
		"object.metadata.ownerReferences":               paths.Object().Metadata().OwnerReferences(),
		"object.metadata.finalizers":                    paths.Object().Metadata().Finalizers(),
		"object.desired":                                paths.Object().Desired(),
		"object.observed":                               paths.Object().Observed(),
		"objectMeta.name":                               paths.ObjectMeta().Name(),
		"objectMeta.generateName":                       paths.ObjectMeta().GenerateName(),
		"objectMeta.namespace":                          paths.ObjectMeta().Namespace(),
		"objectMeta.uid":                                paths.ObjectMeta().UID(),
		"objectMeta.resourceVersion":                    paths.ObjectMeta().ResourceVersion(),
		"objectMeta.generation":                         paths.ObjectMeta().Generation(),
		"objectMeta.createdAt":                          paths.ObjectMeta().CreatedAt(),
		"objectMeta.deletion":                           paths.ObjectMeta().Deletion(),
		"objectMeta.labels":                             paths.ObjectMeta().Labels(),
		"objectMeta.annotations":                        paths.ObjectMeta().Annotations(),
		"objectMeta.ownerReferences":                    paths.ObjectMeta().OwnerReferences(),
		"objectMeta.finalizers":                         paths.ObjectMeta().Finalizers(),
		"ownership.desired":                             paths.Ownership().Desired(),
		"ownership.desired.entries":                     paths.Ownership().DesiredSurface().Entries(),
		"ownership.desired.entries.owner":               paths.Ownership().DesiredSurface().Entry().Owner(),
		"ownership.desired.entries.fields":              paths.Ownership().DesiredSurface().Entry().Fields(),
		"ownership.observed":                            paths.Ownership().Observed(),
		"ownership.observed.entries":                    paths.Ownership().ObservedSurface().Entries(),
		"ownership.observed.entries.owner":              paths.Ownership().ObservedSurface().Entry().Owner(),
		"ownership.observed.entries.fields":             paths.Ownership().ObservedSurface().Entry().Fields(),
		"ownership.metadata":                            paths.Ownership().MetadataPath(),
		"ownership.metadata.labels":                     paths.Ownership().Metadata().Labels(),
		"ownership.metadata.labels.entries":             paths.Ownership().Metadata().LabelsSurface().Entries(),
		"ownership.metadata.labels.entries.owner":       paths.Ownership().Metadata().LabelsSurface().Entry().Owner(),
		"ownership.metadata.labels.entries.fields":      paths.Ownership().Metadata().LabelsSurface().Entry().Fields(),
		"ownership.metadata.annotations":                paths.Ownership().Metadata().Annotations(),
		"ownership.metadata.annotations.entries":        paths.Ownership().Metadata().AnnotationsSurface().Entries(),
		"ownership.metadata.annotations.entries.owner":  paths.Ownership().Metadata().AnnotationsSurface().Entry().Owner(),
		"ownership.metadata.annotations.entries.fields": paths.Ownership().Metadata().AnnotationsSurface().Entry().Fields(),
		"pageMeta.resourceVersion":                      paths.PageMeta().ResourceVersion(),
		"pageMeta.continue":                             paths.PageMeta().Continue(),
		"pageMeta.remainingItemCount":                   paths.PageMeta().RemainingItemCount(),
	}
}
