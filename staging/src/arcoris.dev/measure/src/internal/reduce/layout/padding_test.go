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


package layout

import (
	"testing"
	"unsafe"

	"arcoris.dev/atomicx"
)

func TestPaddedHasAtLeastTwoPadsAroundValue(t *testing.T) {
	var p Padded[uint64]
	if unsafe.Sizeof(p) < uintptr(atomicx.CacheLinePadSize*2)+unsafe.Sizeof(uint64(0)) {
		t.Fatalf("Padded size = %d, too small", unsafe.Sizeof(p))
	}
}
