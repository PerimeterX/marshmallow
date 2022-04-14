// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

import (
	"encoding/json"
	"sync"
	"testing"
)

func BenchmarkMarshmallow(b *testing.B) {
	EnableCache(&sync.Map{})
	s := buildParentStruct()
	data, err := json.Marshal(s)
	if err != nil {
		b.Error("could not marshal data")
		return
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err = Unmarshal(data, &parentStruct{})
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
}

func BenchmarkJSON(b *testing.B) {
	s := buildParentStruct()
	data, err := json.Marshal(s)
	if err != nil {
		b.Error("could not marshal data")
		return
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = json.Unmarshal(data, &parentStruct{})
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
}
