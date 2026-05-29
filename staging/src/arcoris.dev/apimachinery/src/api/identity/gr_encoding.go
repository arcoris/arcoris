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

// MarshalText returns the canonical group/resource text after validation.
func (gr GroupResource) MarshalText() ([]byte, error) {
	return marshalText(gr.String(), gr.Validate)
}

// UnmarshalText parses a canonical group/resource text value.
func (gr *GroupResource) UnmarshalText(data []byte) error {
	if gr == nil {
		return nilReceiver(identityNameGroupResource)
	}

	parsed, err := ParseGroupResource(string(data))
	if err != nil {
		return err
	}

	*gr = parsed
	return nil
}

// MarshalJSON returns the canonical group/resource identity as a JSON string.
func (gr GroupResource) MarshalJSON() ([]byte, error) {
	return marshalJSONString(gr.String(), gr.Validate)
}

// UnmarshalJSON parses a canonical group/resource identity from a JSON string.
func (gr *GroupResource) UnmarshalJSON(data []byte) error {
	if gr == nil {
		return nilReceiver(identityNameGroupResource)
	}

	value, err := unmarshalJSONString(identityNameGroupResource, data)
	if err != nil {
		return err
	}

	return gr.UnmarshalText([]byte(value))
}
