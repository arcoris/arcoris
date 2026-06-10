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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"testing"
)

func TestCompareRefScalar(t *testing.T) {
	resolver := testResolver{
		"example.Name": types.Define("example.Name", types.String()),
	}

	got, err := Compare(value.StringValue("old"), value.StringValue("new"), types.Ref("example.Name").Descriptor(), Options{Resolver: resolver})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.Root()))
}

func TestCompareRefObject(t *testing.T) {
	resolver := testResolver{
		"example.Spec": types.Define(
			"example.Spec",
			types.Object(types.Field("image").String().Optional()),
		),
	}

	got, err := Compare(valueObject("image", "v1"), valueObject("image", "v2"), types.Ref("example.Spec").Descriptor(), Options{Resolver: resolver})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("image")))
}

func TestCompareRefMissingResolver(t *testing.T) {
	_, err := Compare(value.StringValue("old"), value.StringValue("new"), types.Ref("example.Name").Descriptor(), Options{})

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
	requireErrorPath(t, err, "$")
}

func TestCompareRefUnresolved(t *testing.T) {
	_, err := Compare(
		value.StringValue("old"),
		value.StringValue("new"),
		types.Ref("example.Missing").Descriptor(),
		Options{Resolver: testResolver{}},
	)

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
	requireErrorPath(t, err, "$")
}

func TestCompareRefCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.Ref("example.A")),
	}

	_, err := Compare(value.StringValue("old"), value.StringValue("new"), types.Ref("example.A").Descriptor(), Options{Resolver: resolver})

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$")
}

func TestCompareRefDirectCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.A")),
	}

	_, err := Compare(value.StringValue("old"), value.StringValue("new"), types.Ref("example.A").Descriptor(), Options{Resolver: resolver})

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$")
}

func TestCompareRefIndirectCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.Ref("example.C")),
		"example.C": types.Define("example.C", types.Ref("example.A")),
	}

	_, err := Compare(value.StringValue("old"), value.StringValue("new"), types.Ref("example.A").Descriptor(), Options{Resolver: resolver})

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$")
}

func TestCompareRefMaxDepth(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.String()),
	}

	_, err := Compare(
		value.StringValue("old"),
		value.StringValue("new"),
		types.Ref("example.A").Descriptor(),
		Options{Resolver: resolver, MaxDepth: 1},
	)

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$")
}

func TestCompareAtomicListOfRefs(t *testing.T) {
	path := rootField("args")
	resolver := testResolver{
		"example.Name": types.Define("example.Name", types.String()),
	}
	descriptor := types.ListOf(types.Ref("example.Name")).Atomic().Descriptor()

	got, err := CompareAt(
		path,
		value.MustListValue(value.StringValue("old")),
		value.MustListValue(value.StringValue("new")),
		descriptor,
		Options{Resolver: resolver},
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
}
