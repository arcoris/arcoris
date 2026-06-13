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

package objectmemorystore

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestListMatchesResourceAndAllNamespacesScope(t *testing.T) {
	request := objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    objectstore.AllNamespaces(),
	}

	if !listMatches(request, listKey("workers", "alpha", "one")) {
		t.Fatalf("all namespaces scope did not match alpha worker")
	}
	if !listMatches(request, listKey("workers", "beta", "two")) {
		t.Fatalf("all namespaces scope did not match beta worker")
	}
	if listMatches(request, listKey("jobs", "alpha", "skip")) {
		t.Fatalf("request matched a different resource")
	}
}

func TestListMatchesNamespaceScope(t *testing.T) {
	scope, err := objectstore.InNamespace("alpha")
	requireNoError(t, err)
	request := objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    scope,
	}

	if !listMatches(request, listKey("workers", "alpha", "one")) {
		t.Fatalf("namespace scope did not match same namespace")
	}
	if listMatches(request, listKey("workers", "beta", "two")) {
		t.Fatalf("namespace scope matched a different namespace")
	}
	if listMatches(request, listKey("jobs", "alpha", "skip")) {
		t.Fatalf("namespace scope matched a different resource")
	}
}
