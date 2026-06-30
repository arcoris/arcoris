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

package objectcontroller

import (
	"context"
	"fmt"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
	"arcoris.dev/snapshot"
)

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		Scope: objectstore.AllNamespaces(),
	}
}

func testSnapshot(t testing.TB, revision objectstore.Revision) snapshot.Snapshot[objectstore.Revision, objectcache.View] {
	t.Helper()

	cache, err := objectcache.New(testCollection())
	requireNoError(t, err)
	read, err := storewatchapi.NewCollectionRead(testCollection(), objectstore.ListResult{Revision: revision})
	requireNoError(t, err)
	requireNoError(t, cache.Replace(context.Background(), read))
	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	return snap
}

func testItem(id int) objectworkqueue.Item {
	return objectworkqueue.Item{Key: objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{
			Namespace: "default",
			Name:      metaidentity.Name(fmt.Sprintf("worker-%d", id)),
		},
	)}
}
