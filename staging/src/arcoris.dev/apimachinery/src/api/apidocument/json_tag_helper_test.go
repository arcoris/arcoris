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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

const (
	jsonTagName           = "json"
	jsonTagIgnored        = "-"
	jsonTagNoExplicitName = ""
	jsonTagEmptyOption    = ""
	jsonTagSeparator      = ","

	jsonTagOptionInline    = "inline"
	jsonTagOptionOmitEmpty = "omitempty"
	jsonTagOptionOmitZero  = "omitzero"
)

const (
	typeMetaGoFieldAPIVersion = "APIVersion"
	typeMetaGoFieldKind       = "Kind"

	objectGoFieldTypeMeta   = "TypeMeta"
	objectGoFieldObjectMeta = "ObjectMeta"
	objectGoFieldDesired    = "Desired"
	objectGoFieldObserved   = "Observed"

	objectMetaGoFieldName            = "Name"
	objectMetaGoFieldGenerateName    = "GenerateName"
	objectMetaGoFieldNamespace       = "Namespace"
	objectMetaGoFieldUID             = "UID"
	objectMetaGoFieldResourceVersion = "ResourceVersion"
	objectMetaGoFieldGeneration      = "Generation"
	objectMetaGoFieldCreatedAt       = "CreatedAt"
	objectMetaGoFieldDeletion        = "Deletion"
	objectMetaGoFieldLabels          = "Labels"
	objectMetaGoFieldAnnotations     = "Annotations"
	objectMetaGoFieldOwnerReferences = "OwnerReferences"
	objectMetaGoFieldFinalizers      = "Finalizers"
)

// parsedJSONTag is a test-only view of one Go json struct tag.
type parsedJSONTag struct {
	// name is the explicit tag name before the first comma.
	name string

	// options stores comma-separated tag options.
	options map[string]bool

	// ignored reports json:"-".
	ignored bool
}

// parseJSONTag parses one raw json struct tag for drift-prevention tests.
func parseJSONTag(raw string) parsedJSONTag {
	if raw == jsonTagIgnored {
		return parsedJSONTag{ignored: true}
	}

	parts := strings.Split(raw, jsonTagSeparator)
	tag := parsedJSONTag{
		name:    parts[0],
		options: map[string]bool{},
	}
	for _, option := range parts[1:] {
		if option != jsonTagEmptyOption {
			tag.options[option] = true
		}
	}

	return tag
}

// derefType dereferences pointers so helpers can inspect pointer embeddings.
func derefType(goType reflect.Type) reflect.Type {
	for goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}

	return goType
}

// assertJSONTagName asserts that a Go field tag matches a canonical field name.
func assertJSONTagName(
	t *testing.T,
	goType reflect.Type,
	fieldName string,
	want apidocument.FieldName,
) {
	t.Helper()

	assertJSONTagNameString(t, goType, fieldName, want.String())
}

// assertJSONTagNameString asserts that a Go field tag has exactly want as name.
func assertJSONTagNameString(
	t *testing.T,
	goType reflect.Type,
	fieldName string,
	want string,
) {
	t.Helper()

	field := requireStructField(t, goType, fieldName)
	tag := parseJSONTag(field.Tag.Get(jsonTagName))
	if tag.ignored {
		t.Fatalf("%s.%s json tag is ignored", derefType(goType), fieldName)
	}
	if tag.name != want {
		t.Fatalf("%s.%s json tag name = %q; want %q", derefType(goType), fieldName, tag.name, want)
	}
}

// assertJSONTagHasOptions asserts that all wanted tag options are present.
func assertJSONTagHasOptions(
	t *testing.T,
	goType reflect.Type,
	fieldName string,
	wantOptions ...string,
) {
	t.Helper()

	field := requireStructField(t, goType, fieldName)
	tag := parseJSONTag(field.Tag.Get(jsonTagName))
	for _, option := range wantOptions {
		if !tag.options[option] {
			t.Fatalf("%s.%s json tag missing option %q", derefType(goType), fieldName, option)
		}
	}
}

// assertJSONTagLacksOptions asserts that disallowed tag options are absent.
func assertJSONTagLacksOptions(
	t *testing.T,
	goType reflect.Type,
	fieldName string,
	disallowedOptions ...string,
) {
	t.Helper()

	field := requireStructField(t, goType, fieldName)
	tag := parseJSONTag(field.Tag.Get(jsonTagName))
	for _, option := range disallowedOptions {
		if tag.options[option] {
			t.Fatalf("%s.%s json tag unexpectedly has option %q", derefType(goType), fieldName, option)
		}
	}
}

// collectFlattenedJSONFieldNames collects exported document fields in wire order.
func collectFlattenedJSONFieldNames(t *testing.T, goType reflect.Type) []string {
	t.Helper()

	goType = derefType(goType)
	if goType.Kind() != reflect.Struct {
		t.Fatalf("type %s is %s; want struct", goType, goType.Kind())
	}

	var names []string
	for i := 0; i < goType.NumField(); i++ {
		field := goType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		tag := parseJSONTag(field.Tag.Get(jsonTagName))
		if tag.ignored {
			continue
		}
		if field.Anonymous && tag.name == jsonTagNoExplicitName && tag.options[jsonTagOptionInline] {
			names = append(names, collectFlattenedJSONFieldNames(t, field.Type)...)
			continue
		}
		if tag.name == jsonTagNoExplicitName {
			t.Fatalf("%s.%s lacks explicit json field name", goType, field.Name)
		}

		names = append(names, tag.name)
	}

	return names
}

// requireStructField returns a named field or fails the current test.
func requireStructField(t *testing.T, goType reflect.Type, fieldName string) reflect.StructField {
	t.Helper()

	goType = derefType(goType)
	field, ok := goType.FieldByName(fieldName)
	if !ok {
		t.Fatalf("%s has no field %s", goType, fieldName)
	}

	return field
}
