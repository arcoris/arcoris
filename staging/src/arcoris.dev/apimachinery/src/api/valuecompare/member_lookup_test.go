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

package valuecompare

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestMembersByNameBuildsLookup(t *testing.T) {
	got := membersByName([]value.Member{
		value.ObjectMember("env", value.StringValue("prod")),
		value.ObjectMember("tier", value.StringValue("frontend")),
	})

	env, _ := got["env"].String()
	tier, _ := got["tier"].String()
	if env != "prod" || tier != "frontend" {
		t.Fatalf("membersByName() = %#v", got)
	}
}
