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

package meta

import (
	"errors"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/finalizer"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/meta/owner"
	"arcoris.dev/apimachinery/api/meta/stamp"
)

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

func validOwnerReference() owner.Reference {
	return owner.Reference{
		Object: metaidentity.ObjectIdentityReference{
			APIVersion: apiidentity.GroupVersion{Group: "control.arcoris.dev", Version: "v1"},
			Kind:       "Worker",
			Namespace:  "system",
			Name:       "owner",
			UID:        "owner-uid",
		},
		Controller: true,
	}
}

func validObjectMeta() ObjectMeta {
	return ObjectMeta{
		Name:            "worker",
		NamePrefix:      "worker-",
		Namespace:       "system",
		UID:             "uid-1",
		ResourceVersion: "rv-1",
		Generation:      2,
		CreatedAt:       stamp.NewTimestamp(testTime()),
		Labels:          labels.Set{"role": "worker"},
		Annotations:     annotations.Set{"note": "human readable"},
		OwnerReferences: owner.List{validOwnerReference()},
		Finalizers:      finalizer.Set{"cleanup"},
	}
}
