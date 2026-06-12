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

package objectapply

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestApplySameIdentity(t *testing.T) {
	_, err := Apply(testRequest(), Options{})

	requireNoError(t, err)
}

func TestApplyDifferentGVKRejected(t *testing.T) {
	req := testRequest()
	req.Applied.TypeMeta.Kind = apiidentity.Kind("Other")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestApplyDifferentGroupRejected(t *testing.T) {
	req := testRequest()
	req.Applied.TypeMeta = meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
		Group:   "other.arcoris.dev",
		Version: "v1",
		Kind:    "Worker",
	})

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestApplyDifferentVersionRejected(t *testing.T) {
	req := testRequest()
	req.Applied.TypeMeta = testTypeMeta("v2")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrVersionMismatch)
}

func TestApplyDifferentNameRejected(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Name = metaidentity.Name("other")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestApplyDifferentNamespaceRejected(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Namespace = metaidentity.Namespace("other")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestApplyDifferentUIDRejected(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.UID = metaidentity.UID("uid-2")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestApplyAllowsAppliedUIDOmitted(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.UID = ""

	_, err := Apply(req, Options{})

	requireNoError(t, err)
}

func TestApplyAllowsAppliedUIDMatchingLive(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.UID = req.Live.ObjectMeta.UID

	_, err := Apply(req, Options{})

	requireNoError(t, err)
}
