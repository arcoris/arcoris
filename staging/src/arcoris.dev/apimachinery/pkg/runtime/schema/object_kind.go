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

package schema

// ObjectKind is the minimal schema-level contract for values that carry a kind.
//
// The interface deliberately does not mention runtime.Object, metadata, codecs,
// or serializers. Higher layers can embed or adapt this contract without
// creating dependency cycles back into schema.
type ObjectKind interface {
	// SetGroupVersionKind stores the object's schema identity.
	SetGroupVersionKind(kind GroupVersionKind)

	// GroupVersionKind returns the object's current schema identity.
	GroupVersionKind() GroupVersionKind
}

// EmptyObjectKind is a no-op ObjectKind implementation.
//
// It is useful when a higher layer needs to satisfy an ObjectKind-returning
// contract but has no mutable schema identity to store at this layer.
type EmptyObjectKind struct{}

// SetGroupVersionKind intentionally ignores the provided schema identity.
//
// This behavior is useful for immutable or schema-less values that still need
// to satisfy an ObjectKind contract in tests or adapters.
func (EmptyObjectKind) SetGroupVersionKind(GroupVersionKind) {}

// GroupVersionKind always returns the zero schema identity.
//
// The zero value is a deliberate "no kind stored" signal and is not validated
// as a complete identity.
func (EmptyObjectKind) GroupVersionKind() GroupVersionKind {
	return GroupVersionKind{}
}

// ObjectKindHolder is a small mutable ObjectKind implementation.
//
// The holder is intentionally tiny so tests and future runtime layers can store
// a GroupVersionKind without pulling runtime object machinery into schema.
type ObjectKindHolder struct {
	kind GroupVersionKind
}

// NewObjectKindHolder returns a holder initialized with the provided kind.
//
// The constructor does not validate kind because the holder is storage only; the
// caller decides where schema identity validation occurs.
func NewObjectKindHolder(kind GroupVersionKind) *ObjectKindHolder {
	return &ObjectKindHolder{kind: kind}
}

// SetGroupVersionKind stores the provided schema identity.
//
// A nil receiver is ignored so optional embedded holders can be used safely by
// future higher-level layers.
func (h *ObjectKindHolder) SetGroupVersionKind(kind GroupVersionKind) {
	if h == nil {
		return
	}
	h.kind = kind
}

// GroupVersionKind returns the stored schema identity.
//
// A nil receiver returns the zero identity, matching EmptyObjectKind behavior.
func (h *ObjectKindHolder) GroupVersionKind() GroupVersionKind {
	if h == nil {
		return GroupVersionKind{}
	}
	return h.kind
}
