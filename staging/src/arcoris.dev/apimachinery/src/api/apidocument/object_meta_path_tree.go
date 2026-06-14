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

package apidocument

// ObjectMetaPaths groups paths for object metadata documents.
type ObjectMetaPaths struct {
	// base is the document path root for this ObjectMeta view.
	base Path
}

// Name returns the object metadata name path.
func (p ObjectMetaPaths) Name() Path {
	return p.base.Child(Fields().ObjectMeta().Name())
}

// GenerateName returns the object metadata generateName path.
func (p ObjectMetaPaths) GenerateName() Path {
	return p.base.Child(Fields().ObjectMeta().GenerateName())
}

// Namespace returns the object metadata namespace path.
func (p ObjectMetaPaths) Namespace() Path {
	return p.base.Child(Fields().ObjectMeta().Namespace())
}

// UID returns the object metadata uid path.
func (p ObjectMetaPaths) UID() Path {
	return p.base.Child(Fields().ObjectMeta().UID())
}

// ResourceVersion returns the object metadata resourceVersion path.
func (p ObjectMetaPaths) ResourceVersion() Path {
	return p.base.Child(Fields().ObjectMeta().ResourceVersion())
}

// Generation returns the object metadata generation path.
func (p ObjectMetaPaths) Generation() Path {
	return p.base.Child(Fields().ObjectMeta().Generation())
}

// CreatedAt returns the object metadata createdAt path.
func (p ObjectMetaPaths) CreatedAt() Path {
	return p.base.Child(Fields().ObjectMeta().CreatedAt())
}

// Deletion returns the object metadata deletion path.
func (p ObjectMetaPaths) Deletion() Path {
	return p.base.Child(Fields().ObjectMeta().Deletion())
}

// Labels returns the object metadata labels path.
func (p ObjectMetaPaths) Labels() Path {
	return p.base.Child(Fields().ObjectMeta().Labels())
}

// Annotations returns the object metadata annotations path.
func (p ObjectMetaPaths) Annotations() Path {
	return p.base.Child(Fields().ObjectMeta().Annotations())
}

// OwnerReferences returns the object metadata ownerReferences path.
func (p ObjectMetaPaths) OwnerReferences() Path {
	return p.base.Child(Fields().ObjectMeta().OwnerReferences())
}

// Finalizers returns the object metadata finalizers path.
func (p ObjectMetaPaths) Finalizers() Path {
	return p.base.Child(Fields().ObjectMeta().Finalizers())
}
