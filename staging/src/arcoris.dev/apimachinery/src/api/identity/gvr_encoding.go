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

// MarshalText returns the canonical group/version/resource text after validation.
func (gvr GroupVersionResource) MarshalText() ([]byte, error) {
	return marshalText(gvr.String(), gvr.Validate)
}

// UnmarshalText parses a canonical group/version/resource text value.
func (gvr *GroupVersionResource) UnmarshalText(data []byte) error {
	if gvr == nil {
		return nilReceiver(identityNameGroupVersionResource)
	}

	parsed, err := ParseGroupVersionResource(string(data))
	if err != nil {
		return err
	}

	*gvr = parsed
	return nil
}

// MarshalJSON returns the canonical group/version/resource identity as a JSON string.
func (gvr GroupVersionResource) MarshalJSON() ([]byte, error) {
	return marshalJSONString(gvr.String(), gvr.Validate)
}

// UnmarshalJSON parses a canonical group/version/resource identity from a JSON string.
func (gvr *GroupVersionResource) UnmarshalJSON(data []byte) error {
	if gvr == nil {
		return nilReceiver(identityNameGroupVersionResource)
	}

	value, err := unmarshalJSONString(identityNameGroupVersionResource, data)
	if err != nil {
		return err
	}

	return gvr.UnmarshalText([]byte(value))
}
