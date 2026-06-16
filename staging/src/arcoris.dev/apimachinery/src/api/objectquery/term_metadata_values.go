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

import "sort"

// canonicalStrings returns a detached sorted set of unique metadata values.
func canonicalStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	out := append([]string(nil), values...)
	sort.Strings(out)
	next := out[:0]
	for _, value := range out {
		if len(next) == 0 || next[len(next)-1] != value {
			next = append(next, value)
		}
	}

	return append([]string(nil), next...)
}

// metadataDomainName returns stable diagnostic text for metadata domains.
func metadataDomainName(domain metadataDomain) string {
	switch domain {
	case metadataLabels:
		return "label"
	case metadataAnnotations:
		return "annotation"
	default:
		return "metadata"
	}
}
