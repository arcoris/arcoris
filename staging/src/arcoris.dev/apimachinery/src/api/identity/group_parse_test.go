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

package identity

import (
	"strings"
	"testing"
)

func TestParseGroupCanonicalValues(t *testing.T) {
	label63 := strings.Repeat("a", 63)
	total253 := strings.Repeat("a", 63) + dnsLabelSeparator +
		strings.Repeat("b", 63) + dnsLabelSeparator +
		strings.Repeat("c", 63) + dnsLabelSeparator +
		strings.Repeat("d", 61)

	requireCanonicalSet(t, []string{
		"",
		"control.arcoris.dev",
		"runtime.arcoris.dev",
		"a.b",
		label63 + dnsLabelSeparator + "dev",
		total253,
	}, ParseGroup)
}

func TestParseGroupRejectsNonCanonicalValues(t *testing.T) {
	label64 := strings.Repeat("a", 64)
	total254 := strings.Repeat("a", 63) + dnsLabelSeparator +
		strings.Repeat("b", 63) + dnsLabelSeparator +
		strings.Repeat("c", 63) + dnsLabelSeparator +
		strings.Repeat("d", 62)

	requireRejectedSet(t, []string{
		"apps",
		"batch",
		"arcoris",
		"Control.arcoris.dev",
		"control..arcoris.dev",
		"control_arcoris.dev",
		"control/arcoris/dev",
		" control.arcoris.dev",
		"control.arcoris.dev ",
		label64 + dnsLabelSeparator + "dev",
		total254,
	}, ParseGroup)
}
