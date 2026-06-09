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

package objectlifecycle

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestApplyMapsStoreStaleRevision(t *testing.T) {
	store := staleUpdateStore{state: committedStateForFakeStore()}
	executor, err := NewExecutor(
		WithStore(store),
		WithResourceResolver(testCatalog(t)),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	requireNoError(t, err)

	_, err = executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrStaleRevision, ReasonStaleRevision)
	requireErrorIs(t, err, objectstore.ErrStaleRevision)
}

func committedStateForFakeStore() objectstore.State {
	return objectstore.State{
		Object: testObject(1, "api:v1"),
		Ownership: objectownership.Document{
			Version: objectownership.VersionV1,
			Desired: objectownership.Surface{
				Entries: []objectownership.Entry{
					{Owner: owner("creator"), Fields: []objectownership.Path{"$.image"}},
				},
			},
		},
		Revision: 1,
	}
}

type staleUpdateStore struct {
	state objectstore.State
}

func (s staleUpdateStore) Get(context.Context, objectstore.Key) (objectstore.State, bool, error) {
	return s.state, true, nil
}

func (s staleUpdateStore) Create(context.Context, objectstore.Key, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, objectstore.ErrAlreadyExists
}

func (s staleUpdateStore) Update(context.Context, objectstore.Key, objectstore.Revision, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, objectstore.ErrStaleRevision
}

func (s staleUpdateStore) Delete(context.Context, objectstore.Key, objectstore.Revision) (objectstore.State, error) {
	return objectstore.State{}, objectstore.ErrStaleRevision
}
