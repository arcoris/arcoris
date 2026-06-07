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

package healthregistry

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestBuilderRegisterEmptyBatchIsNoop(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	if err := builder.Register(health.TargetReady); err != nil {
		t.Fatalf("Register(empty) = %v, want nil", err)
	}

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}
	if !registry.Empty() {
		t.Fatal("empty registration should not create registry entries")
	}
}

func TestBuilderRegisterRejectsExistingDuplicateWithoutPartialCommit(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	mustRegister(t, builder, health.TargetReady, checker("storage"))

	err := builder.Register(health.TargetReady, checker("storage"), checker("queue"))
	if !errors.Is(err, ErrDuplicateCheck) {
		t.Fatalf("Register() = %v, want ErrDuplicateCheck", err)
	}

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}
	if registry.Has(health.TargetReady, "queue") {
		t.Fatal("failed registration batch committed queue")
	}
}
