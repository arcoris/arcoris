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

package objectstore

import (
	"context"
	"errors"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func validResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

func validObjectName() metaidentity.ObjectName {
	return metaidentity.ObjectName{Namespace: "system", Name: "main"}
}

func validKey() Key {
	return MustKey(validResource(), validObjectName())
}

func validState() State {
	return State{
		Object: object.New[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{
				Name:      "main",
				Namespace: "system",
			},
			value.StringValue("desired"),
		),
		Ownership: objectownership.Document{Version: objectownership.VersionV1},
	}
}

func validCommittedState() State {
	state := validState()
	state.Revision = 1

	return state
}

func ownershipWithEntry() objectownership.Document {
	return objectownership.Document{
		Version: objectownership.VersionV1,
		Desired: objectownership.Surface{
			Entries: []objectownership.Entry{
				{
					Owner:  "manager",
					Fields: []objectownership.Path{"$.spec"},
				},
			},
		},
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

type fakeStore struct{}

func (fakeStore) Get(context.Context, Key) (State, bool, error) {
	return State{}, false, nil
}

func (fakeStore) Create(context.Context, Key, State) (State, error) {
	return State{}, nil
}

func (fakeStore) Update(context.Context, Key, Revision, State) (State, error) {
	return State{}, nil
}

func (fakeStore) Delete(context.Context, Key, Revision) (State, error) {
	return State{}, nil
}
