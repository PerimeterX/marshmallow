# Marshmallow

Package marshmallow provides a simple API to perform flexible and performant JSON unmarshalling. Unlike other packages,
marshmallow supports unmarshalling of some known and some unknown fields with zero performance overhead nor extra coding
needed. While unmarshalling, marshmallow allows fully retaining the original data and access it via a typed struct and a
dynamic map.

- [Install](#install)
- [Usage](#usage)
    * [Where Does Marshmallow Shine](#where-does-marshmallow-shine)
- [Alternatives and Performance Benchmark](#alternatives-and-performance-benchmark)
- [API](#api)
    * [Unmarshal](#unmarshal)
    * [UnmarshalFromJSONMap](#unmarshalfromjsonmap)
    * [API Options](#api-options)
    * [Caching](#caching)

## Install

```sh
go get -u github.com/perimeterx/marshmallow
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/perimeterx/marshmallow"
	"sync"
)

func main() {
	marshmallow.EnableCache(&sync.Map{}) // this is used to boost performance, read more below
	v := struct {
		Foo string `json:"foo"`
		Boo []int  `json:"boo"`
	}{}
	result, err := marshmallow.Unmarshal([]byte(`{"foo":"bar","boo":[1,2,3],"goo":12.6}`), &v)
	fmt.Printf("v=%+v, result=%+v, err=%v", v, result, err)
	// Output: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar goo:12.6], err=<nil>
}
```

#### Where Does Marshmallow Shine

Marshmallow is best suited for use cases where you are interested in all the input data, but you have predetermined
information only about a subset of it. For instance, if you plan to reference two specific fields from the data, then
iterate all the data and apply some generic logic. How does it look with the native library:

```go
func isAllowedToDrive(data []byte) (bool, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal(data, &result)
	if err != nil {
		return false, err
	}

	age, ok := result["age"]
	if !ok {
		return false, nil
	}
	a, ok := age.(float64)
	if !ok {
		return false, nil
	}
	if a < 17 {
		return false, nil
	}

	hasDriversLicense, ok := result["has_drivers_license"]
	if !ok {
		return false, nil
	}
	h, ok := hasDriversLicense.(bool)
	if !ok {
		return false, nil
	}
	if !h {
		return false, nil
	}

	for key := range result {
		if strings.Contains(key, "prior_conviction") {
			return false, nil
		}
	}

	return true, nil
}
```

And with marshmallow:

```go
func isAllowedToDrive(data []byte) (bool, error) {
	v := struct {
		Age               int  `json:"age"`
		HasDriversLicense bool `json:"has_drivers_license"`
	}{}
	result, err := marshmallow.Unmarshal(data, &v)
	if err != nil {
		return false, err
	}

	if v.Age < 17 || !v.HasDriversLicense {
		return false, nil
	}

	for key := range result {
		if strings.Contains(key, "prior_conviction") {
			return false, nil
		}
	}

	return true, nil
}
```

There can be two main reasons to have an interest in all the data. First is when you eventually plan to write or pipe
the input data, and you don't want to lose any of it. Second is if you plan to perform any kind of dynamic read of
the data - this includes iterating it, reading calculated or configured field names, and others.

## Alternatives and Performance Benchmark

[Full Benchmark](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go)

Other solutions available for this kind of use case, each solution is explained and documented in the link below.

|Benchmark|(1)|(2)|(3)|(4)|
|--|--|--|--|--|
|[unmarshall twice](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L40)|228693|5164 ns/op|1640 B/op|51 allocs/op|
|[raw map](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L66)|232236|5116 ns/op|2296 B/op|53 allocs/op|
|[go codec](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L121)|388442|3077 ns/op|2512 B/op|37 allocs/op|
|[marshmallow](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L16)|626168|1853 ns/op|608 B/op|18 allocs/op|
|[marshmallow without populating struct](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L162)|678616|1751 ns/op|608 B/op|18 allocs/op|

**Marshmallow provides the best performance (up to X3 faster) while not requiring any extra coding.**
In fact, marshmallow performs as fast as normal `json.Unmarshal` call, however, it populates both the map and the
struct.

|Benchmark|(1)|(2)|(3)|(4)|
|--|--|--|--|--|
|[marshmallow](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L16)|626168|1853 ns/op|608 B/op|18 allocs/op|
|[native library](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L143)|652106|1845 ns/op|304 B/op|11 allocs/op|
|[marshmallow without populating struct](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L162)|678616|1751 ns/op|608 B/op|18 allocs/op|

## API

Marshmallow exposes two main API functions

#### Unmarshal

`marshmallow.Unmarshal(data []byte, v interface{}, options ...UnmarshalOption) (map[string]interface{}, error)`

Unmarshal parses the JSON-encoded object in data and stores the values in the struct pointed to by v and in the returned
map. If v is nil or not a pointer to a struct, Unmarshal returns an ErrInvalidValue. If data is not a valid JSON or not
a JSON object Unmarshal returns an ErrInvalidInput.

Unmarshal follows the rules of json.Unmarshal with the following exceptions:

- All input fields are stored in the resulting map, including fields that do not exist in the struct pointed by v.
- Unmarshal only operates on JSON object inputs. It will reject all other types of input by returning ErrInvalidInput.
- Unmarshal only operates on struct values. It will reject all other types of v by returning ErrInvalidValue.
- Unmarshal supports three types of Mode values. Each mode is documented below.

#### UnmarshalFromJSONMap

`marshmallow.UnmarshalFromJSONMap(data map[string]interface{}, v interface{}, options ...UnmarshalOption) (map[string]interface{}, error)`

UnmarshalFromJSONMap parses the JSON map data and stores the values in the struct pointed to by v and in the returned
map. If v is nil or not a pointer to a struct, UnmarshalFromJSONMap returns an ErrInvalidValue.

UnmarshalFromJSONMap follows the rules of json.Unmarshal with the following exceptions:

- All input fields are stored in the resulting map, including fields that do not exist in the struct pointed by v.
- UnmarshalFromJSONMap receive a JSON map instead of raw bytes. The given input map is assumed to be a JSON map, meaning
  it should only contain the following types: `bool`, `string`, `float64`, `[]interface`, and `map[string]interface{}`.
  Other types will cause decoding to return unexpected results.
- UnmarshalFromJSONMap only operates on struct values. It will reject all other types of v by returning ErrInvalidValue.
- UnmarshalFromJSONMap supports three types of Mode values. Each mode is documented below.

**UnmarshalerFromJSONMap** is the interface implemented by types that can unmarshal a JSON description of themselves. In
case you want to implement custom unmarshalling, json.Unmarshaler only supports receiving the data as []byte. However,
while unmarshalling from JSON map, the data is not available as a raw []byte and converting to it will significantly
hurt performance. Thus, if you wish to implement a custom unmarshalling on a type that is being unmarshalled from a JSON
map, you need to implement UnmarshalerFromJSONMap interface.

#### API Options

- `marshmallow.WithMode(mode Mode)` sets the unmarshalling mode:
    - **ModeFailOnFirstError** is the default mode. It makes unmarshalling terminate immediately on any kind of error.
      This error will then be returned.
    - **ModeAllowMultipleErrors** mode makes unmarshalling keep decoding even if errors are encountered. In case of such
      error, the erroneous value will be omitted from the result. Eventually, all errors will all be returned, alongside
      the partial result.
    - **ModeFailOverToOriginalValue** mode makes unmarshalling keep decoding even if errors are encountered. In case of
      such error, the original external value be placed in the result data, even though it does not meet the schematic
      requirements. Eventually, all errors will be returned, alongside the full result. Note that the result map
      will contain values that do not match the struct schema.
- `marshmallow.WithSkipPopulateStruct(skipPopulateStruct bool)` sets the skipPopulateStruct option. Skipping populate
  struct is set to false by default. If you do not intend to use the struct value once unmarshalling is finished, set
  this option to true to boost performance. This would mean the struct fields will not be set with values, but rather it
  will only be used as the target schema when populating the result map.

#### Caching

`marshmallow.EnableCache` enables unmarshalling cache. It allows reuse of refection information about types needed to
perform the unmarshalling. A use of such cache can boost up unmarshalling by x1.4. Check out
[benchmark_test.go](benchmark_test.go) for an example.

`EnableCache` is not thread safe! Do not use it while performing unmarshalling, or it will cause an unsafe race condition.
Typically, `EnableCache` should be called once when the process boots.

Caching is disabled by default. The use of this function allows enabling it and controlling the behavior of the cache.
Typically, the use of `sync.Map` should be good enough. The caching mechanism stores a single `map` per struct type. If
you plan to unmarshal a huge amount of distinct struct it may get to consume a lot of resources, in which case you have
the control to choose the caching implementation you like and its setup.
