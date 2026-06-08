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

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

// DecodeBinding binds one normalized decode key to one configured codec entry.
type DecodeBinding struct {
	// ContentType is the normalized input content key.
	ContentType ContentType

	// Target is the API document model being decoded.
	Target codec.Target

	// Transport is the required byte or stream decode capability.
	Transport Transport

	// EntryID is the exact configured codec candidate to select.
	EntryID codecregistry.EntryID
}
