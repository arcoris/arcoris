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

package object

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/meta"
)

func TestListValidateMeta(t *testing.T) {
	list := NewList(
		validListTypeMeta(),
		validPageMeta(),
		[]uninspectedPayload{{Value: "item"}},
	)

	requireNoError(t, list.ValidateMeta())
}

func TestListValidateMetaRejectsInvalidMetadata(t *testing.T) {
	tests := []struct {
		name   string
		list   List[int]
		target error
		path   string
		reason ErrorReason
	}{
		{
			name: "type meta",
			list: List[int]{
				TypeMeta: meta.TypeMeta{Kind: "WorkerList"},
				PageMeta: validPageMeta(),
			},
			target: meta.ErrInvalidTypeMeta,
			path:   "list.typeMeta",
			reason: ErrorReasonInvalidTypeMeta,
		},
		{
			name: "list meta",
			list: List[int]{
				TypeMeta: validListTypeMeta(),
				PageMeta: meta.PageMeta{ContinueToken: "bad token"},
			},
			target: meta.ErrInvalidPageMeta,
			path:   "list.metadata",
			reason: ErrorReasonInvalidPageMeta,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.list.ValidateMeta()
			requireErrorIs(t, err, ErrInvalidList)
			requireErrorIs(t, err, tt.target)

			var objectErr *Error
			if !errors.As(err, &objectErr) {
				t.Fatalf("errors.As(%T) = false", objectErr)
			}
			if objectErr.Path != tt.path {
				t.Fatalf("Path = %q, want %q", objectErr.Path, tt.path)
			}
			if objectErr.Reason != tt.reason {
				t.Fatalf("Reason = %q", objectErr.Reason)
			}
			if objectErr.Cause == nil {
				t.Fatal("Cause = nil")
			}
			if objectErr.Detail == "" {
				t.Fatal("Detail = empty")
			}
		})
	}
}

func TestListValidateMetaDoesNotInspectItems(t *testing.T) {
	called := false
	list := NewList(
		validListTypeMeta(),
		validPageMeta(),
		[]payloadWithValidate{{Called: &called}},
	)

	requireNoError(t, list.ValidateMeta())
	if called {
		t.Fatal("ValidateMeta called item Validate")
	}
}
