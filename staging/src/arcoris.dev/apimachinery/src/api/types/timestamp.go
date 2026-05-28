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

// TimestampType builds descriptors for absolute points in time.
//
// TimestampType records an absolute point-in-time descriptor. It does not
// choose RFC3339, Unix time, or another encoding; codecs own concrete
// representations.
type TimestampType struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
	// payload stores the exact timestamp constraints under construction.
	payload timestampPayload
}

// Timestamp returns a timestamp descriptor builder.
//
// Typical reusable declaration:
//
//	observedAtType := Timestamp().Nullable()
func Timestamp() TimestampType {
	return TimestampType{header: newHeader(TypeTimestamp)}
}

// Nullable returns a timestamp descriptor that admits null values.
func (t TimestampType) Nullable() TimestampType {
	t.header = t.header.withNullable()
	return t
}

// Type returns a detached Type descriptor.
func (t TimestampType) Type() Type {
	out := typeFromHeader(t.header)
	out.timestamp = cloneTimestampPayload(t.payload)
	return out
}

// typeExpr marks TimestampType as a sealed TypeExpr implementation.
func (t TimestampType) typeExpr() {}
