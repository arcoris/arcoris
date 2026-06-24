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

package jsonconfig

import "testing"

func TestDefaultEncodeOutputConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeOutputConfig()

	if config.Layout != LayoutCompact {
		t.Fatalf("layout = %d; want compact", config.Layout)
	}
	if config.Indent != defaultIndent {
		t.Fatalf("indent = %q; want %q", config.Indent, defaultIndent)
	}
	if config.FinalNewline != FinalNewlineOmit {
		t.Fatalf("final newline = %d; want omit", config.FinalNewline)
	}
}

func TestResolveEncodeOutputConfig(t *testing.T) {
	t.Parallel()

	config := EncodeOutputConfig{}
	resolveEncodeOutputConfig(&config)

	if config.Layout != LayoutCompact {
		t.Fatalf("layout = %d; want compact", config.Layout)
	}
	if config.Indent != defaultIndent {
		t.Fatalf("indent = %q; want %q", config.Indent, defaultIndent)
	}
	if config.FinalNewline != FinalNewlineOmit {
		t.Fatalf("final newline = %d; want omit", config.FinalNewline)
	}
}

func TestValidateEncodeOutputConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config EncodeOutputConfig
		path   string
	}{
		"layout": {
			config: EncodeOutputConfig{Layout: LayoutMode(99), Indent: defaultIndent, FinalNewline: FinalNewlineOmit},
			path:   "encode.output.layout",
		},
		"final newline": {
			config: EncodeOutputConfig{Layout: LayoutCompact, Indent: defaultIndent, FinalNewline: FinalNewlineMode(99)},
			path:   "encode.output.final_newline",
		},
		"pretty empty indent": {
			config: EncodeOutputConfig{Layout: LayoutPretty, Indent: "", FinalNewline: FinalNewlineOmit},
			path:   "encode.output.indent",
		},
		"control indent": {
			config: EncodeOutputConfig{Layout: LayoutPretty, Indent: "\n", FinalNewline: FinalNewlineOmit},
			path:   "encode.output.indent",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeOutputConfig(testCase.config)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}

func TestHasInvalidIndentRune(t *testing.T) {
	t.Parallel()

	if hasInvalidIndentRune(" \t") {
		t.Fatalf("space and tab indent reported invalid")
	}
	if !hasInvalidIndentRune(".") {
		t.Fatalf("non-indent rune reported valid")
	}
}
