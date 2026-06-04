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

package codecjson

import "testing"

func TestResolvedConfigZeroValue(t *testing.T) {
	t.Parallel()

	var config resolvedConfig

	if config.decode != (resolvedDecodeConfig{}) {
		t.Fatalf("decode = %#v; want zero", config.decode)
	}
	if config.encode != (resolvedEncodeConfig{}) {
		t.Fatalf("encode = %#v; want zero", config.encode)
	}
}
