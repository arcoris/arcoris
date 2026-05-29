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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestObjectLike(t *testing.T) {
	ok, detail := objectLike(
		objectType(),
		nil,
		make(map[types.TypeName]bool),
		"desired",
	)
	requireEqual(t, ok, true)
	requireEqual(t, detail, "")
}

func TestObjectLikeThroughResolver(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.Object"): types.Define("control.arcoris.dev.Object", types.Object()),
	}

	ok, detail := objectLike(
		refType("control.arcoris.dev.Object"),
		resolver,
		make(map[types.TypeName]bool),
		"desired",
	)
	requireEqual(t, ok, true)
	requireEqual(t, detail, "")
}

func TestObjectLikeThroughNestedResolverRefs(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.A"): types.Define(
			"control.arcoris.dev.A",
			types.Ref("control.arcoris.dev.B"),
		),
		types.TypeName("control.arcoris.dev.B"): types.Define(
			"control.arcoris.dev.B",
			types.Object(),
		),
	}

	ok, detail := objectLike(
		refType("control.arcoris.dev.A"),
		resolver,
		make(map[types.TypeName]bool),
		"desired",
	)
	requireEqual(t, ok, true)
	requireEqual(t, detail, "")
}

func TestObjectLikeRejectsNonObjectRoots(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.Text"): types.Define("control.arcoris.dev.Text", types.String()),
	}

	cases := []struct {
		name   string
		typ    types.Type
		r      types.Resolver
		detail string
	}{
		{name: "scalar", typ: stringType(), detail: "object or reference to object"},
		{name: "unresolved ref", typ: refType("control.arcoris.dev.Object"), detail: "requires a resolver"},
		{name: "ref not found", typ: refType("control.arcoris.dev.Missing"), r: resolver, detail: "was not found"},
		{name: "ref scalar", typ: refType("control.arcoris.dev.Text"), r: resolver, detail: "object or reference to object"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, detail := objectLike(
				tc.typ,
				tc.r,
				make(map[types.TypeName]bool),
				"desired",
			)
			requireEqual(t, ok, false)
			if !strings.Contains(detail, tc.detail) {
				t.Fatalf("detail = %q, want to contain %q", detail, tc.detail)
			}
		})
	}
}

func TestObjectLikeRejectsRecursiveReferences(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.Loop"): types.Define(
			"control.arcoris.dev.Loop",
			types.Ref("control.arcoris.dev.Loop"),
		),
	}

	ok, detail := objectLike(
		refType("control.arcoris.dev.Loop"),
		resolver,
		make(map[types.TypeName]bool),
		"desired",
	)
	requireEqual(t, ok, false)

	if !strings.Contains(detail, "recursive") {
		t.Fatalf("detail = %q, want recursive reference detail", detail)
	}
}

func TestRequireObjectLikeReturnsVersionError(t *testing.T) {
	err := requireObjectLike(
		stringType(),
		nil,
		"definition.versions[v1].desired",
		ErrorReasonDesiredNotObject,
		"desired",
	)

	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].desired",
		ErrorReasonDesiredNotObject,
	)
}
