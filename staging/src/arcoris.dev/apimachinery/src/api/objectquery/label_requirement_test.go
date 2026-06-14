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

package objectquery

import "testing"

func TestLabelRequirementZeroValue(t *testing.T) {
	var req LabelRequirement
	if req.req.op != 0 || req.req.key != "" || req.req.values != nil {
		t.Fatalf("zero label requirement = %#v; want empty", req)
	}
}

func mustLabelExists(t *testing.T, key string) LabelRequirement {
	t.Helper()
	req, err := LabelExists(key)
	requireNoError(t, err)

	return req
}

func mustLabelDoesNotExist(t *testing.T, key string) LabelRequirement {
	t.Helper()
	req, err := LabelDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustLabelEquals(t *testing.T, key string, value string) LabelRequirement {
	t.Helper()
	req, err := LabelEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustLabelNotEquals(t *testing.T, key string, value string) LabelRequirement {
	t.Helper()
	req, err := LabelNotEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustLabelIn(t *testing.T, key string, values ...string) LabelRequirement {
	t.Helper()
	req, err := LabelIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelNotIn(t *testing.T, key string, values ...string) LabelRequirement {
	t.Helper()
	req, err := LabelNotIn(key, values...)
	requireNoError(t, err)

	return req
}
