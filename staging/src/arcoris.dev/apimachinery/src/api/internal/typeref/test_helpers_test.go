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

package typeref

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

type resolverFunc func(types.TypeName) (types.Definition, bool)

func (f resolverFunc) Resolve(name types.TypeName) (types.Definition, bool) {
	return f(name)
}

func stringResolver(name types.TypeName) (types.Definition, bool) {
	return types.Define(name, types.String()), true
}

func exampleResolver() types.Resolver {
	return resolverFunc(func(name types.TypeName) (types.Definition, bool) {
		if name == "example.Name" {
			return types.Define("example.Name", types.String()), true
		}

		return types.Definition{}, false
	})
}

func chainResolver() types.Resolver {
	return resolverFunc(func(name types.TypeName) (types.Definition, bool) {
		switch name {
		case "example.Name":
			return types.Define("example.Name", types.Ref("example.Text")), true
		case "example.Text":
			return types.Define("example.Text", types.String()), true
		default:
			return types.Definition{}, false
		}
	})
}

func cycleResolver() types.Resolver {
	return resolverFunc(func(name types.TypeName) (types.Definition, bool) {
		switch name {
		case "example.A":
			return types.Define("example.A", types.Ref("example.B")), true
		case "example.B":
			return types.Define("example.B", types.Ref("example.A")), true
		default:
			return types.Definition{}, false
		}
	})
}

func rootPath() fieldpath.Path {
	return fieldpath.RootPath()
}

func refType(name types.TypeName) types.Descriptor {
	return types.Ref(name).Descriptor()
}

func requireFailureKind(t *testing.T, err error, want FailureKind) {
	t.Helper()

	refError, ok := AsError(err)
	if !ok {
		t.Fatalf("error = %v; want typeref.Error", err)
	}
	if refError.Kind != want {
		t.Fatalf("error kind = %s; want %s", refError.Kind, want)
	}
}
