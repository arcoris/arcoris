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

package objectstorewatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestNewRejectsNilBackend(t *testing.T) {
	store, err := New(nil)

	if store != nil {
		t.Fatalf("store = %#v; want nil", store)
	}
	requireErrorIs(t, err, ErrNilBackend)
}

func TestNewWithDefaultsSucceeds(t *testing.T) {
	store := testRuntimeStore(t)

	if store == nil {
		t.Fatalf("store is nil")
	}
}

func TestStoreInterfaces(t *testing.T) {
	var _ objectstore.Store = (*Store)(nil)
	var _ storewatchapi.CollectionLister = (*Store)(nil)
	var _ storewatchapi.ListerWatcher = (*Store)(nil)
	var _ storewatchapi.Store = (*Store)(nil)
	var _ storewatchapi.CapableStore = (*Store)(nil)
	var _ objectwatch.Source = (*Store)(nil)
	var _ objectwatch.CapabilityReporter = (*Store)(nil)
}
