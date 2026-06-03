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

package valuecompare

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

type testResolver map[types.TypeName]types.TypeDefinition

func (r testResolver) ResolveType(name types.TypeName) (types.TypeDefinition, bool) {
	definition, ok := r[name]
	return definition, ok
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireErrorNotIs(t *testing.T, err, target error) {
	t.Helper()
	if errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = true", err, target)
	}
}

func requireErrorPath(t *testing.T, err error, want string) {
	t.Helper()

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false for %v", err)
	}
	if got.Path != want {
		t.Fatalf("error path = %q, want %q", got.Path, want)
	}
}

func requireErrorReason(t *testing.T, err error, want ErrorReason) {
	t.Helper()

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false for %v", err)
	}
	if got.Reason != want {
		t.Fatalf("error reason = %q, want %q", got.Reason, want)
	}
}

func requireErrorDetailContains(t *testing.T, err error, want string) {
	t.Helper()

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false for %v", err)
	}
	if !strings.Contains(got.Detail, want) {
		t.Fatalf("error detail = %q, want substring %q", got.Detail, want)
	}
}

func requireResult(
	t *testing.T,
	got Result,
	added []fieldpath.Path,
	removed []fieldpath.Path,
	modified []fieldpath.Path,
) {
	t.Helper()

	requireDisjointResult(t, got)
	requireSet(t, "added", got.Added, added...)
	requireSet(t, "removed", got.Removed, removed...)
	requireSet(t, "modified", got.Modified, modified...)
}

func requireDisjointResult(t *testing.T, result Result) {
	t.Helper()

	if overlap := result.Added.Intersection(result.Removed); !overlap.IsEmpty() {
		t.Fatalf("added/removed overlap = %s", overlap)
	}
	if overlap := result.Added.Intersection(result.Modified); !overlap.IsEmpty() {
		t.Fatalf("added/modified overlap = %s", overlap)
	}
	if overlap := result.Removed.Intersection(result.Modified); !overlap.IsEmpty() {
		t.Fatalf("removed/modified overlap = %s", overlap)
	}

	union := result.Added.Union(result.Removed).Union(result.Modified)
	if !result.Changed().Equal(union) {
		t.Fatalf("changed = %s, want union %s", result.Changed(), union)
	}
}

func requireSet(t *testing.T, name string, got fieldpath.Set, want ...fieldpath.Path) {
	t.Helper()

	expected := fieldpath.MustSet(want...)
	if !got.Equal(expected) {
		t.Fatalf("%s = %s, want %s", name, got, expected)
	}
}

func requireNoChangedPathContaining(t *testing.T, result Result, fragment string) {
	t.Helper()

	for _, p := range result.Changed().Paths() {
		if strings.Contains(p.String(), fragment) {
			t.Fatalf("changed path %s contains %q", p, fragment)
		}
	}
}

func rootField(names ...string) fieldpath.Path {
	path := fieldpath.RootPath()
	for _, name := range names {
		path = path.Field(name)
	}

	return path
}

func paths(values ...fieldpath.Path) []fieldpath.Path {
	return values
}

func mustDecimal(t *testing.T, text string) value.Value {
	t.Helper()

	decimal, err := value.NewDecimal(text)
	if err != nil {
		t.Fatalf("NewDecimal(%q) error = %v", text, err)
	}

	return value.DecimalValue(decimal)
}

func mustFloat(t *testing.T, f float64) value.Value {
	t.Helper()

	v, err := value.FloatValue(f)
	if err != nil {
		t.Fatalf("FloatValue(%v) error = %v", f, err)
	}

	return v
}

func readySelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("type", fieldpath.StringLiteral("Ready")),
	)
}

func degradedSelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("type", fieldpath.StringLiteral("Degraded")),
	)
}

func progressingSelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("type", fieldpath.StringLiteral("Progressing")),
	)
}

func routeSelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("port", fieldpath.Uint64Literal(443)),
		fieldpath.NewSelectorEntry("host", fieldpath.StringLiteral("api.example.com")),
	)
}

func conditionDescriptor() types.Type {
	return conditionExpr().Type()
}

func conditionExpr() types.ObjectType {
	return types.Object(
		types.Field("type").String().Required(),
		types.Field("status").String().Required(),
	)
}

func conditionsDescriptor() types.Type {
	return types.ListOf(
		types.Object(
			types.Field("type").String().Required(),
			types.Field("status").String().Required(),
		),
	).Map("type").Type()
}

func conditionValue(conditionType string, status string) value.Value {
	return value.MustObjectValue(
		value.ObjectMember("type", value.StringValue(conditionType)),
		value.ObjectMember("status", value.StringValue(status)),
	)
}

func valueObject(fields ...string) value.Value {
	members := make([]value.Member, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		members = append(members, value.ObjectMember(fields[i], value.StringValue(fields[i+1])))
	}

	return value.MustObjectValue(members...)
}

func typesObject(fields ...string) types.Type {
	exprs := make([]types.FieldExpr, 0, len(fields))
	for _, name := range fields {
		exprs = append(exprs, types.Field(name).String().Optional())
	}

	return types.Object(exprs...).Type()
}

func imageContainer(image string) value.Value {
	return value.MustObjectValue(
		value.ObjectMember("name", value.StringValue("main")),
		value.ObjectMember("image", value.StringValue(image)),
	)
}
