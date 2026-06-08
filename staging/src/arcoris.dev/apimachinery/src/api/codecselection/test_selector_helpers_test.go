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
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

func testSelector(t *testing.T, config Config) Selector {
	t.Helper()

	selector, err := New(config)
	requireNoError(t, err)

	return selector
}

func testSelectorForDecode(
	t *testing.T,
	registration codecregistry.Registration,
	target codec.Target,
	transport Transport,
) Selector {
	t.Helper()

	return testSelector(t, Config{
		Registry: testRegistry(t, registration),
		DecodeBindings: []DecodeBinding{{
			ContentType: testContentType(codec.MediaTypeJSON),
			Target:      target,
			Transport:   transport,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
	})
}

func testSelectorForEncode(
	t *testing.T,
	registration codecregistry.Registration,
	target codec.Target,
	transport Transport,
) Selector {
	t.Helper()

	return testSelector(t, Config{
		Registry: testRegistry(t, registration),
		EncodeBindings: []EncodeBinding{{
			ContentType: testContentType(codec.MediaTypeJSON),
			Target:      target,
			Transport:   transport,
			EntryID:     codecregistry.MustEntryID("json.public"),
		}},
	})
}
