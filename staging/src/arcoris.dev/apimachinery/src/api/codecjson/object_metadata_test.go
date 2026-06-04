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

package codecjson

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

// TestNodeToOptionalObjectMetaAbsent covers omitted metadata.
func TestNodeToOptionalObjectMetaAbsent(t *testing.T) {
	got, err := nodeToOptionalObjectMeta(rootPath(), jsonNode{kind: jsonKindObject}, newTestCodec(t).decode)
	requireNoError(t, err)

	if !got.IsZero() {
		t.Fatalf("metadata = %#v; want zero", got)
	}
}

// TestObjectMetaToNodeRoundTripsThroughMetadataDecoder covers metadata delegation.
func TestObjectMetaToNodeRoundTripsThroughMetadataDecoder(t *testing.T) {
	want := testObjectMeta()
	metadataNode, err := objectMetaToNode(rootPath().Member(apidocument.ObjectFieldMetadata.String()), want, newTestCodec(t).decode)
	requireNoError(t, err)

	envelopeNode := jsonNode{
		kind: jsonKindObject,
		members: []jsonMember{
			{name: apidocument.ObjectFieldMetadata.String(), value: metadataNode},
		},
	}
	got, err := nodeToOptionalObjectMeta(rootPath(), envelopeNode, newTestCodec(t).decode)
	requireNoError(t, err)

	if got.Name != want.Name {
		t.Fatalf("name = %q; want %q", got.Name, want.Name)
	}
	if got.Namespace != want.Namespace {
		t.Fatalf("namespace = %q; want %q", got.Namespace, want.Namespace)
	}
}
