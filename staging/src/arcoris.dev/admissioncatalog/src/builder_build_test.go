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

package admissioncatalog

import (
	"errors"
	"testing"
)

func TestBuilderBuildsIndependentCatalogs(t *testing.T) {
	var builder Builder
	if err := builder.DeclareKind(kindDescriptor(testKind)); err != nil {
		t.Fatalf("DeclareKind returned error: %v", err)
	}
	first, err := builder.Build()
	if err != nil {
		t.Fatalf("first Build returned error: %v", err)
	}

	if err := builder.DeclareReason(reasonDescriptor(testReason)); err != nil {
		t.Fatalf("DeclareReason returned error: %v", err)
	}
	second, err := builder.Build()
	if err != nil {
		t.Fatalf("second Build returned error: %v", err)
	}

	if first.HasReason(testReason) {
		t.Fatal("first catalog changed after builder mutation")
	}
	if !second.HasReason(testReason) {
		t.Fatal("second catalog is missing later declaration")
	}
}

func TestBuilderMustBuildPanicsOnInvalidState(t *testing.T) {
	builder := newInitializedBuilder()
	builder.components.declare(componentDescriptor(testComponent, testKind))

	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("MustBuild did not panic")
		}
		err, ok := got.(error)
		if !ok {
			t.Fatalf("panic = %T, want error", got)
		}
		if !errors.Is(err, ErrUnknownComponentKind) {
			t.Fatalf("panic error = %v, want ErrUnknownComponentKind", err)
		}
	}()
	_ = builder.MustBuild()
}
