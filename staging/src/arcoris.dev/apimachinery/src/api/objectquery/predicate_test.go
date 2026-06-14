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

package objectquery

import (
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func testItem(namespace string, name string, labelValues map[string]string, annotationValues map[string]string) objectstore.ListItem {
	labelSet, err := labels.FromStrings(labelValues)
	if err != nil {
		panic(err)
	}
	annotationSet, err := annotations.FromStrings(annotationValues)
	if err != nil {
		panic(err)
	}

	key := objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace(namespace),
			Name:      metaidentity.Name(name),
		},
	)

	return objectstore.ListItem{
		Key: key,
		State: objectstore.State{
			Object: object.New[value.Value, value.Value](
				meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
					Group:   "control.arcoris.dev",
					Version: "v1",
					Kind:    "Worker",
				}),
				meta.ObjectMeta{
					Name:        metaidentity.Name(name),
					Namespace:   metaidentity.Namespace(namespace),
					Labels:      labelSet,
					Annotations: annotationSet,
				},
				value.StringValue(name),
			),
			Revision: 1,
		},
	}
}
