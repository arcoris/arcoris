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

package apidocument_test

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestObjectFieldNames(t *testing.T) {
	assertFieldName(t, "ObjectFieldAPIVersion", apidocument.ObjectFieldAPIVersion, "apiVersion")
	assertFieldName(t, "ObjectFieldKind", apidocument.ObjectFieldKind, "kind")
	assertFieldName(t, "ObjectFieldMetadata", apidocument.ObjectFieldMetadata, "metadata")
	assertFieldName(t, "ObjectFieldDesired", apidocument.ObjectFieldDesired, "desired")
	assertFieldName(t, "ObjectFieldObserved", apidocument.ObjectFieldObserved, "observed")
}
