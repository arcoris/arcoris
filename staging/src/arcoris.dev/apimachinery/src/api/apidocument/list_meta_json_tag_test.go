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
	"slices"
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
	"arcoris.dev/apimachinery/api/meta"
)

func TestListMetaJSONTagsMatchAPIDocument(t *testing.T) {
	goType := reflect.TypeOf(meta.ListMeta{})
	fields := []struct {
		goName string
		api    apidocument.FieldName
	}{
		{goName: listMetaGoFieldResourceVersion, api: apidocument.ListMetaFieldResourceVersion},
		{goName: listMetaGoFieldContinueToken, api: apidocument.ListMetaFieldContinue},
		{goName: listMetaGoFieldRemainingItemCount, api: apidocument.ListMetaFieldRemainingItemCount},
	}

	for _, field := range fields {
		t.Run(field.goName, func(t *testing.T) {
			assertJSONTagName(t, goType, field.goName, field.api)
			assertJSONTagHasOptions(t, goType, field.goName, jsonTagOptionOmitEmpty)
		})
	}
}

func TestListMetaFlattenedJSONFieldsMatchAPIDocument(t *testing.T) {
	goType := reflect.TypeOf(meta.ListMeta{})
	got := collectFlattenedJSONFieldNames(t, goType)
	want := []string{
		apidocument.ListMetaFieldResourceVersion.String(),
		apidocument.ListMetaFieldContinue.String(),
		apidocument.ListMetaFieldRemainingItemCount.String(),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("list metadata fields = %#v; want %#v", got, want)
	}
}
