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

package metagrammar

import "arcoris.dev/apimachinery/api/internal/lexical"

// ValidateDNSLabel validates one lowercase DNS-1123-like metadata label.
func ValidateDNSLabel(s string) *Violation {
	return fromLexical(lexical.ValidateDNS1123Label(s))
}
