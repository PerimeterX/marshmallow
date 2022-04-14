// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

import (
	"reflect"
)

type Cache interface {
	Load(key interface{}) (interface{}, bool)
	Store(key, value interface{})
}

func EnableCache(c Cache) {
	cache = c
}

var cache Cache

func cacheLookup(t reflect.Type) map[string]int {
	if cache == nil {
		return nil
	}
	value, exists := cache.Load(t)
	if !exists {
		return nil
	}
	result, _ := value.(map[string]int)
	return result
}

func cacheStore(t reflect.Type, fields map[string]int) {
	if cache == nil {
		return
	}
	cache.Store(t, fields)
}
