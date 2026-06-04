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
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
	"arcoris.dev/apimachinery/api/meta"
)

func TestTypeMetaJSONTagsMatchAPIDocument(t *testing.T) {
	typ := reflect.TypeOf(meta.TypeMeta{})

	assertJSONTagName(t, typ, typeMetaGoFieldAPIVersion, apidocument.TypeMetaFieldAPIVersion)
	assertJSONTagHasOptions(t, typ, typeMetaGoFieldAPIVersion, jsonTagOptionOmitEmpty, jsonTagOptionOmitZero)

	assertJSONTagName(t, typ, typeMetaGoFieldKind, apidocument.TypeMetaFieldKind)
	assertJSONTagHasOptions(t, typ, typeMetaGoFieldKind, jsonTagOptionOmitEmpty)
	assertJSONTagLacksOptions(t, typ, typeMetaGoFieldKind, jsonTagOptionOmitZero)
}
