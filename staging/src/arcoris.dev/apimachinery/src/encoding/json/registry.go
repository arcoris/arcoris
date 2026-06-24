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

package jsoncodec

import (
	"arcoris.dev/apimachinery/api/codecregistry"
	"arcoris.dev/apimachinery/encoding/json/jsonconfig"
)

// Registration builds an explicit registry declaration for one configured JSON
// codec candidate.
//
// The caller supplies the registry identity because api/codecregistry treats
// EntryID as owner policy, not codec metadata. Registration constructs and
// validates the JSON codec configuration, but leaves EntryID validation and
// registry conflict checks to codecregistry.New.
func Registration(id codecregistry.EntryID, config jsonconfig.Config) (codecregistry.Registration, error) {
	c, err := New(config)
	if err != nil {
		return codecregistry.Registration{}, err
	}

	return codecregistry.Register(id, c), nil
}
