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

package resource

import (
	"testing"

	"arcoris.dev/apimachinery/api/identity"
)

type resolverFunc struct {
	resource        func(identity.GroupResource) (Definition, bool)
	kind            func(identity.GroupKind) (Definition, bool)
	versionResource func(identity.GroupVersionResource) (Definition, VersionDefinition, bool)
	versionKind     func(identity.GroupVersionKind) (Definition, VersionDefinition, bool)
}

func (f resolverFunc) ResolveResource(gr identity.GroupResource) (Definition, bool) {
	return f.resource(gr)
}

func (f resolverFunc) ResolveKind(gk identity.GroupKind) (Definition, bool) {
	return f.kind(gk)
}

func (f resolverFunc) ResolveVersionResource(
	gvr identity.GroupVersionResource,
) (Definition, VersionDefinition, bool) {
	return f.versionResource(gvr)
}

func (f resolverFunc) ResolveVersionKind(
	gvk identity.GroupVersionKind,
) (Definition, VersionDefinition, bool) {
	return f.versionKind(gvk)
}

func TestResolverContract(t *testing.T) {
	def := validDefinition()
	version, ok := def.Version(identity.Version("v1"))
	if !ok {
		t.Fatalf("test definition is missing v1")
	}

	resolver := resolverFunc{
		resource: func(identity.GroupResource) (Definition, bool) {
			return def, true
		},
		kind: func(identity.GroupKind) (Definition, bool) {
			return def, true
		},
		versionResource: func(identity.GroupVersionResource) (Definition, VersionDefinition, bool) {
			return def, version, true
		},
		versionKind: func(identity.GroupVersionKind) (Definition, VersionDefinition, bool) {
			return def, version, true
		},
	}

	var contract Resolver = resolver
	gvr, ok := def.GroupVersionResource("v1")
	if !ok {
		t.Fatalf("test definition is missing v1 GVR")
	}
	gvk, ok := def.GroupVersionKind("v1")
	if !ok {
		t.Fatalf("test definition is missing v1 GVK")
	}

	if _, ok := contract.ResolveResource(def.GroupResource()); !ok {
		t.Fatalf("ResolveResource() ok = false")
	}
	if _, ok := contract.ResolveKind(def.GroupKind()); !ok {
		t.Fatalf("ResolveKind() ok = false")
	}
	if _, _, ok := contract.ResolveVersionResource(gvr); !ok {
		t.Fatalf("ResolveVersionResource() ok = false")
	}
	if _, _, ok := contract.ResolveVersionKind(gvk); !ok {
		t.Fatalf("ResolveVersionKind() ok = false")
	}
}
