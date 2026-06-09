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

func TestTypeNameValidationMatrix(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"", false},
		{"Name", false},
		{"meta.arcoris.dev.Name", true},
		{"meta.arcoris.dev.Namespace", true},
		{"arcoris.schema.Group", true},
		{"example.dev.CronExpression", true},
		{"example.dev.cronExpression", false},
		{"Bad.CronExpression", false},
		{"example.dev..CronExpression", false},
		{"example.dev.Cron-Expression", false},
		{"example-name.CronExpression", true},
		{"example_name.CronExpression", false},
		{"example.dev.Имя", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, TypeName(tt.name).IsValid(), tt.valid)
			_, err := ParseTypeName(tt.name)
			requireEqual(t, err == nil, tt.valid)
		})
	}
}
