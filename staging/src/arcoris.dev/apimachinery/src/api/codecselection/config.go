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

package codecselection

import "arcoris.dev/apimachinery/api/codecregistry"

// Config describes an explicit immutable selector plan to build.
type Config struct {
	// Registry is the already-built configured codec candidate catalog.
	Registry codecregistry.Registry

	// DecodeBindings bind normalized decode keys to exact EntryID values.
	DecodeBindings []DecodeBinding

	// EncodeBindings bind normalized encode keys to exact EntryID values.
	EncodeBindings []EncodeBinding
}
