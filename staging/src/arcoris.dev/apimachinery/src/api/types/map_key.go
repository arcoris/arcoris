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

// MapKeyType identifies the structural key family for map descriptors.
type MapKeyType uint8

const (
	// MapKeyString is the only supported map key type in this design pass.
	MapKeyString MapKeyType = iota
)

// IsValid reports whether k is a supported map key type.
func (k MapKeyType) IsValid() bool {
	return k == MapKeyString
}
