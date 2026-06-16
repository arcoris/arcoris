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

package objectcache

import (
	"testing"

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectquery"
)

func mustNamespaceEquals(t *testing.T, namespace metaidentity.Namespace) objectquery.NamespaceRequirement {
	t.Helper()

	req, err := objectquery.NamespaceEquals(namespace)
	requireNoError(t, err)

	return req
}

func mustNameEquals(t *testing.T, name metaidentity.Name) objectquery.NameRequirement {
	t.Helper()

	req, err := objectquery.NameEquals(name)
	requireNoError(t, err)

	return req
}

func mustLabelExists(t *testing.T, key string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelExists(key)
	requireNoError(t, err)

	return req
}

func mustLabelDoesNotExist(t *testing.T, key string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustLabelEquals(t *testing.T, key string, val string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustLabelNotEquals(t *testing.T, key string, val string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelNotEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustLabelIn(t *testing.T, key string, values ...string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelNotIn(t *testing.T, key string, values ...string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelNotIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelSelector(
	t *testing.T,
	requirements ...objectquery.LabelRequirement,
) objectquery.LabelSelector {
	t.Helper()

	selector, err := objectquery.NewLabelSelector(requirements...)
	requireNoError(t, err)

	return selector
}

func mustAnnotationExists(t *testing.T, key string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationExists(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationDoesNotExist(t *testing.T, key string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationEquals(t *testing.T, key string, val string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotEquals(t *testing.T, key string, val string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationNotEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustAnnotationIn(t *testing.T, key string, values ...string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotIn(t *testing.T, key string, values ...string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationNotIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationSelector(
	t *testing.T,
	requirements ...objectquery.AnnotationRequirement,
) objectquery.AnnotationSelector {
	t.Helper()

	selector, err := objectquery.NewAnnotationSelector(requirements...)
	requireNoError(t, err)

	return selector
}
