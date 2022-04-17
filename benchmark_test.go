// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

import (
	"encoding/json"
	"github.com/ugorji/go/codec"
	"sync"
	"testing"
)

func BenchmarkMarshmallow(b *testing.B) {
	EnableCache(&sync.Map{})
	var v benchmarkParent
	var result map[string]interface{}
	var err error
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v = benchmarkParent{}
		result, err = Unmarshal(benchmarkData, &v)
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
	b.StopTimer()
	validateBenchmarkStruct(b, &v)
	validateBenchmarkTypedMap(b, result)
}

func BenchmarkUnmarshalTwice(b *testing.B) {
	var v benchmarkParent
	var result map[string]interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v = benchmarkParent{}
		err := json.Unmarshal(benchmarkData, &v)
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
		result = make(map[string]interface{})
		err = json.Unmarshal(benchmarkData, &result)
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
	b.StopTimer()
	validateBenchmarkStruct(b, &v)
	validateBenchmarkUntypedMap(b, result)
}

func BenchmarkUnmarshalRawMap(b *testing.B) {
	var v benchmarkParent
	var result map[string]interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data := make(map[string]json.RawMessage)
		err := json.Unmarshal(benchmarkData, &data)
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
		v = benchmarkParent{}
		result = make(map[string]interface{})
		for key, value := range data {
			switch key {
			case "field1":
				err = json.Unmarshal(value, &v.Field1)
				if err != nil {
					b.Error("could not unmarshal data")
					return
				}
				result["field1"] = v.Field1
			case "field2":
				err = json.Unmarshal(value, &v.Field2)
				if err != nil {
					b.Error("could not unmarshal data")
					return
				}
				result["field2"] = v.Field2
			case "field3":
				err = json.Unmarshal(value, &v.Field3)
				if err != nil {
					b.Error("could not unmarshal data")
					return
				}
				result["field3"] = v.Field3
			default:
				var i interface{}
				err = json.Unmarshal(value, &i)
				if err != nil {
					b.Error("could not unmarshal data")
					return
				}
				result[key] = i
			}
		}
	}
	b.StopTimer()
	validateBenchmarkStruct(b, &v)
	validateBenchmarkTypedMap(b, result)
}

func BenchmarkGoCodec(b *testing.B) {
	var v benchmarkParent
	var result map[string]interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v = benchmarkParent{}
		result = make(map[string]interface{})
		result["field1"] = &v.Field1
		result["field2"] = &v.Field2
		result["field3"] = &v.Field3
		err := codec.NewDecoderBytes(benchmarkData, &codec.JsonHandle{}).Decode(&result)
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
	b.StopTimer()
	validateBenchmarkStruct(b, &v)
}

func BenchmarkJSON(b *testing.B) {
	b.ResetTimer()
	var v benchmarkParent
	for n := 0; n < b.N; n++ {
		v = benchmarkParent{}
		err := json.Unmarshal(benchmarkData, &v)
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
	b.StopTimer()
	validateBenchmarkStruct(b, &v)
}

func BenchmarkMarshmallowWithSkipPopulateStruct(b *testing.B) {
	EnableCache(&sync.Map{})
	var v benchmarkParent
	var result map[string]interface{}
	var err error
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v = benchmarkParent{}
		result, err = Unmarshal(benchmarkData, &v, WithSkipPopulateStruct(true))
		if err != nil {
			b.Error("could not unmarshal data")
			return
		}
	}
	b.StopTimer()
	validateBenchmarkTypedMap(b, result)
}

var benchmarkData = []byte(`{"field1":"foo","field2":12,"field3":{"field1":"boo","field2":24},"field4":[1,24,false],"field5":false}`)

func validateBenchmarkStruct(b *testing.B, result *benchmarkParent) {
	if result.Field1 == "foo" && result.Field2 == 12 && result.Field3.Field1 == "boo" && result.Field3.Field2 == 24 {
		return
	}
	b.Error("invalid struct data")
}

func validateBenchmarkTypedMap(b *testing.B, m map[string]interface{}) {
	if m["field1"] == "foo" &&
		m["field2"] == 12 &&
		m["field3"].(*benchmarkChild).Field1 == "boo" &&
		m["field3"].(*benchmarkChild).Field2 == 24 &&
		m["field4"].([]interface{})[0] == float64(1) &&
		m["field4"].([]interface{})[1] == float64(24) &&
		m["field4"].([]interface{})[2] == false &&
		m["field5"] == false {
		return
	}
	b.Error("invalid map data")
}

func validateBenchmarkUntypedMap(b *testing.B, m map[string]interface{}) {
	if m["field1"] == "foo" &&
		m["field2"] == float64(12) &&
		m["field3"].(map[string]interface{})["field1"] == "boo" &&
		m["field3"].(map[string]interface{})["field2"] == float64(24) &&
		m["field4"].([]interface{})[0] == float64(1) &&
		m["field4"].([]interface{})[1] == float64(24) &&
		m["field4"].([]interface{})[2] == false &&
		m["field5"] == false {
		return
	}
	b.Error("invalid map data")
}

type benchmarkParent struct {
	Field1 string          `json:"field1"`
	Field2 int             `json:"field2"`
	Field3 *benchmarkChild `json:"field3"`
}

type benchmarkChild struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}
