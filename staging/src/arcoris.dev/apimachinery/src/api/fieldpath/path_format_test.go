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
		{"root", Root(), "$"},
		{"field", Root().Field(testField("spec")).Field(testField("replicas")), "$.spec.replicas"},
		{"key app", Root().Field(testField("metadata")).Field(testField("labels")).Key(testKey("app")), `$.metadata.labels["app"]`},
		{
			"key dotted",
			Root().Field(testField("metadata")).Field(testField("labels")).Key(testKey("app.kubernetes.io/name")),
			`$.metadata.labels["app.kubernetes.io/name"]`,
		},
		{"index", Root().Field(testField("containers")).Index(0).Field(testField("image")), "$.containers[0].image"},
		{
			"selector one key",
			Root().
				Field("conditions").
				Select(MustSelector(testSelectorEntry("type", StringLiteral("Ready")))).
				Field("status"),
			`$.conditions[{"type":"Ready"}].status`,
		},
		{
			"selector multi key",
			Root().
				Field("ports").
				Select(
					MustSelector(
						testSelectorEntry("protocol", StringLiteral("TCP")),
						testSelectorEntry("name", StringLiteral("http")),
					),
				).
				Field("port"),
			`$.ports[{"name":"http","protocol":"TCP"}].port`,
		},
		{
			"selector route",
			Root().
				Field("routes").
				Select(
					MustSelector(
						testSelectorEntry("port", Uint64Literal(443)),
						testSelectorEntry("host", StringLiteral("api.example.com")),
					),
				).
				Field("backend"),
			`$.routes[{"host":"api.example.com","port":443}].backend`,
		},
		{"quoted field", Root().Field(testField("api-version")), `$."api-version"`},
		{
			"quoted field and selector",
			Root().
				Field("x-y-z").
				Select(MustSelector(testSelectorEntry("name", StringLiteral(`a"b`)))),
			`$."x-y-z"[{"name":"a\"b"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.path.String(), tt.want)
		})
	}
}
