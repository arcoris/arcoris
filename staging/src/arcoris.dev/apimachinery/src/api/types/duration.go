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

// DurationType builds descriptors for elapsed intervals.
//
// DurationType records an elapsed-interval descriptor. It does not choose
// ISO-8601, nanoseconds, Go duration text, or another encoding.
type DurationType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact duration constraints under construction.
	payload durationPayload
}

// Duration returns a duration descriptor builder.
//
// Typical reusable declaration:
//
//	timeoutType := Duration().Nullable()
func Duration() DurationType {
	return DurationType{header: newHeader(TypeDuration)}
}

// Nullable returns a duration descriptor that admits null values.
func (t DurationType) Nullable() DurationType {
	t.header = t.header.withNullable()
	return t
}

// Type returns a detached Type descriptor.
func (t DurationType) Type() Type {
	out := typeFromHeader(t.header)
	out.duration = cloneDurationPayload(t.payload)
	return out
}

// typeExpr marks DurationType as a sealed TypeExpr implementation.
func (t DurationType) typeExpr() {}
