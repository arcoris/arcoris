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

import "testing"

func TestUint8ValidateRejectsInvalidRules(t *testing.T) {
	tests := []Type{
		Uint8().Range(10, 1).Type(),
		Uint8().Min(1).Enum(0).Type(),
		Uint8().Max(1).Enum(2).Type(),
		Uint8().Enum(1, 1).Type(),
	}
	for _, typ := range tests {
		requireErrorIs(t, ValidateType(typ, nil), ErrInvalidType)
	}
}
