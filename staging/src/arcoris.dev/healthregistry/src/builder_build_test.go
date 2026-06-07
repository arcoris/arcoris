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
	"testing"

	"arcoris.dev/health"
)

func TestBuilderBuildReturnsImmutableRegistrySnapshot(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	mustRegister(t, builder, health.TargetReady, checker("storage"))

	registry, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	mustRegister(t, builder, health.TargetReady, checker("later"))
	if registry.Has(health.TargetReady, "later") {
		t.Fatal("built registry changed after later builder mutation")
	}
}
