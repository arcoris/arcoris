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

import "unicode"

const defaultIndent = "  "

// EncodeOutputConfig controls JSON output layout and document framing.
type EncodeOutputConfig struct {
	// Layout controls compact or pretty output.
	Layout LayoutMode

	// Indent is used when Layout is LayoutPretty.
	Indent string

	// FinalNewline controls whether the encoded document ends with '\n'.
	FinalNewline FinalNewlineMode
}

// defaultEncodeOutputConfig returns compact output without a final newline.
func defaultEncodeOutputConfig() EncodeOutputConfig {
	return EncodeOutputConfig{
		Layout:       LayoutCompact,
		Indent:       defaultIndent,
		FinalNewline: FinalNewlineOmit,
	}
}

// resolveEncodeOutputConfig applies output defaults in place.
func resolveEncodeOutputConfig(config *EncodeOutputConfig) {
	defaultLayout := config.Layout == LayoutDefault
	if config.Layout == LayoutDefault {
		config.Layout = LayoutCompact
	}
	if config.Indent == "" && defaultLayout {
		config.Indent = defaultIndent
	}
	if config.FinalNewline == FinalNewlineDefault {
		config.FinalNewline = FinalNewlineOmit
	}
}

// validateEncodeOutputConfig checks layout, indentation, and framing policy.
func validateEncodeOutputConfig(config EncodeOutputConfig) error {
	switch {
	case !isKnownLayoutMode(config.Layout):
		return invalidConfig("encode.output.layout", "unknown layout mode %d", config.Layout)
	case !isKnownFinalNewlineMode(config.FinalNewline):
		return invalidConfig("encode.output.final_newline", "unknown final newline mode %d", config.FinalNewline)
	case config.Layout == LayoutPretty && config.Indent == "":
		return invalidConfig("encode.output.indent", "must be non-empty for pretty layout")
	case config.Indent != "" && hasInvalidIndentRune(config.Indent):
		return invalidConfig("encode.output.indent", "must contain only spaces or tabs")
	default:
		return nil
	}
}

// hasInvalidIndentRune reports whether indent contains non-space indentation.
func hasInvalidIndentRune(indent string) bool {
	for _, r := range indent {
		if r == ' ' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) || r == '\n' || r == '\r' {
			return true
		}
		return true
	}

	return false
}
