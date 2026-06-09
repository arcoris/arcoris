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

package typeref

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestResolveFinalReturnsNonRefDescriptor(t *testing.T) {
	descriptor, err := New(nil, 64).ResolveFinal(
		rootPath(),
		types.String().Descriptor(),
		0,
	)
	if err != nil {
		t.Fatalf("ResolveFinal() error = %v", err)
	}
	if descriptor.Code() != types.DescriptorString {
		t.Fatalf("ResolveFinal() descriptor = %s; want %s", descriptor.Code(), types.DescriptorString)
	}
}

func TestResolveFinalResolvesReferenceChain(t *testing.T) {
	descriptor, err := New(chainResolver(), 64).ResolveFinal(
		rootPath(),
		refType("example.Name"),
		0,
	)
	if err != nil {
		t.Fatalf("ResolveFinal() error = %v", err)
	}
	if descriptor.Code() != types.DescriptorString {
		t.Fatalf("ResolveFinal() descriptor = %s; want %s", descriptor.Code(), types.DescriptorString)
	}
}

func TestResolveFinalRejectsIndirectCycle(t *testing.T) {
	_, err := New(cycleResolver(), 64).ResolveFinal(
		rootPath(),
		refType("example.A"),
		0,
	)
	requireFailureKind(t, err, FailureReferenceCycle)
}

func TestResolveFinalRejectsMaxDepth(t *testing.T) {
	_, err := New(chainResolver(), 1).ResolveFinal(
		rootPath(),
		refType("example.Name"),
		0,
	)
	requireFailureKind(t, err, FailureReferenceCycle)
}
