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

package objectstore_test

import (
	"fmt"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

func ExampleKey() {
	key := objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{Namespace: "system", Name: "main"},
	)

	fmt.Println(key)

	// Output:
	// control.arcoris.dev/v1:workers/system/main
}
