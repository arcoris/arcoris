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
	"fmt"

	"arcoris.dev/apimachinery/api/objectmemorystore"
	"arcoris.dev/apimachinery/api/resourcecatalog"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func ExampleExecutor_Create() {
	store, _ := objectmemorystore.New()
	catalog := testCatalogForExample()
	executor, _ := NewExecutor(
		WithStore(store),
		WithResourceResolver(catalog),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)

	result, _ := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	fmt.Println(result.Effect, result.State.Revision.IsValid())
	// Output: created true
}

func ExampleExecutor_Apply() {
	store, _ := objectmemorystore.New()
	catalog := testCatalogForExample()
	executor, _ := NewExecutor(
		WithStore(store),
		WithResourceResolver(catalog),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	_, _ = executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	result, _ := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator")},
	)

	fmt.Println(result.Effect)
	// Output: updated
}

func ExampleExecutor_Delete() {
	store, _ := objectmemorystore.New()
	catalog := testCatalogForExample()
	executor, _ := NewExecutor(
		WithStore(store),
		WithResourceResolver(catalog),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	created, _ := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	result, _ := executor.Delete(
		context.Background(),
		DeleteRequest{Resource: testGVR(), Object: testName(1), Expected: created.State.Revision},
	)

	fmt.Println(result.Effect)
	// Output: deleted
}

func testCatalogForExample() ResourceResolver {
	catalog := resourcecatalog.New(nil)
	if err := catalog.Register(testDefinition()); err != nil {
		panic(err)
	}
	return catalog
}
