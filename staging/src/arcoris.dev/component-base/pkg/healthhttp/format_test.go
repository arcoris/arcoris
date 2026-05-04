/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package healthhttp

import "testing"

func TestFormatString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{name: "text", format: FormatText, want: "text"},
		{name: "json", format: FormatJSON, want: "json"},
		{name: "invalid", format: Format(99), want: "invalid"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.format.String(); got != test.want {
				t.Fatalf("String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestFormatIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format Format
		want   bool
	}{
		{name: "text", format: FormatText, want: true},
		{name: "json", format: FormatJSON, want: true},
		{name: "invalid", format: Format(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.format.IsValid(); got != test.want {
				t.Fatalf("IsValid() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestFormatContentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{name: "text", format: FormatText, want: contentTypeText},
		{name: "json", format: FormatJSON, want: contentTypeJSON},
		{name: "invalid", format: Format(99), want: ""},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.format.contentType(); got != test.want {
				t.Fatalf("contentType() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestFormatZeroValueIsText(t *testing.T) {
	t.Parallel()

	var format Format
	if format != FormatText {
		t.Fatalf("zero Format = %s, want %s", format, FormatText)
	}
	if !format.IsValid() {
		t.Fatal("zero Format should be valid")
	}
	if got := format.contentType(); got != contentTypeText {
		t.Fatalf("zero Format contentType() = %q, want %q", got, contentTypeText)
	}
}

func TestValidateFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format Format
		want   bool
	}{
		{name: "text", format: FormatText, want: true},
		{name: "json", format: FormatJSON, want: true},
		{name: "invalid", format: Format(99), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validateFormat(test.format)
			if got := err == nil; got != test.want {
				t.Fatalf("validateFormat(%s) ok = %v, want %v; err=%v", test.format, got, test.want, err)
			}
		})
	}
}
