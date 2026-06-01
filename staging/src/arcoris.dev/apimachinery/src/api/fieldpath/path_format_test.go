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

package fieldpath

import "testing"

func TestPathStringExamples(t *testing.T) {
	tests := []struct {
		name string
		path Path
		want string
	}{
		{"root", RootPath(), "$"},
		{"field", RootPath().Field("spec").Field("replicas"), "$.spec.replicas"},
		{"key app", RootPath().Field("metadata").Field("labels").Key("app"), `$.metadata.labels["app"]`},
		{"key dotted", RootPath().Field("metadata").Field("labels").Key("app.kubernetes.io/name"), `$.metadata.labels["app.kubernetes.io/name"]`},
		{"index", RootPath().Field("containers").Index(0).Field("image"), "$.containers[0].image"},
		{"selector one key", RootPath().Field("conditions").Select(MustSelector(NewSelectorEntry("type", StringLiteral("Ready")))).Field("status"), `$.conditions[{"type":"Ready"}].status`},
		{"selector multi key", RootPath().Field("ports").Select(MustSelector(NewSelectorEntry("protocol", StringLiteral("TCP")), NewSelectorEntry("name", StringLiteral("http")))).Field("port"), `$.ports[{"name":"http","protocol":"TCP"}].port`},
		{"selector route", RootPath().Field("routes").Select(MustSelector(NewSelectorEntry("port", Uint64Literal(443)), NewSelectorEntry("host", StringLiteral("api.example.com")))).Field("backend"), `$.routes[{"host":"api.example.com","port":443}].backend`},
		{"quoted field", RootPath().Field("api-version"), `$."api-version"`},
		{"quoted field and selector", RootPath().Field("x-y-z").Select(MustSelector(NewSelectorEntry("name", StringLiteral(`a"b`)))), `$."x-y-z"[{"name":"a\"b"}]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.path.String(), tt.want)
		})
	}
}
