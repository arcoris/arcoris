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

package types

// cloneDefinition returns def with descriptor payload storage detached.
//
// Definition is a small value type, but Descriptor can contain slice-backed
// payloads. Catalogs and resolvers use this helper before publishing or
// returning definitions so callers cannot mutate stored descriptor state.
func cloneDefinition(def Definition) Definition {
	def.descriptor = cloneDescriptor(def.descriptor)

	return def
}
