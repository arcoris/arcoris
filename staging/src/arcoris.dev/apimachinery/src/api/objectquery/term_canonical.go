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

import "strings"

// canonicalKey encodes a leaf term into a deterministic internal sort key. It
// is not a public query syntax and may change with the private representation.
func (t term) canonicalKey() string {
	switch t.kind {
	case termResource:
		return "resource=" + t.resource.String()
	case termNamespace:
		return "namespace=" + t.namespace.String()
	case termName:
		return "name=" + t.name.String()
	case termObject:
		return "object=" + t.namespace.String() + "/" + t.name.String()
	case termKey:
		return "key=" + t.key.String()
	case termMetadata:
		return metadataDomainName(t.metadataDomain) + "." +
			t.metadataKey + "." + t.operator.String() + "=" +
			strings.Join(t.stringValues, "|")
	case termField:
		values := make([]string, 0, len(t.values))
		for _, literal := range t.values {
			values = append(values, canonicalValueKey(literal))
		}
		return "field." + t.fieldRef.String() + "." + t.operator.String() +
			"=" + strings.Join(values, "|")
	default:
		return "unknown"
	}
}
