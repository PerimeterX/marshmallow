# Marshmallow

![marshmallow](https://raw.githubusercontent.com/PerimeterX/marshmallow/assets/logo.png)

Package marshmallow provides a simple API to perform flexible and performant JSON unmarshalling. Unlike other packages,
marshmallow supports unmarshalling of some known and some unknown fields with zero performance overhead nor extra coding
needed. While unmarshalling, marshmallow allows fully retaining the original data and access it via a typed struct and a
dynamic map.

- [Marshmallow](#marshmallow)
    * [Install](#install)
    * [Usage](#usage)
    * [Performance Benchmark And Alternatives](#performance-benchmark-and-alternatives)
    * [When Should I Use Marshmallow](#when-should-i-use-marshmallow)
    * [API](#api)

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
)

func main() {
	marshmallow.EnableCache() // this is used to boost performance, read more below
	v := struct {
		Foo string `json:"foo"`
		Boo []int  `json:"boo"`
	}{}
	result, err := marshmallow.Unmarshal([]byte(`{"foo":"bar","boo":[1,2,3],"goo":12.6}`), &v)
	fmt.Printf("v=%+v, result=%+v, err=%v", v, result, err)
	// Output: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar goo:12.6], err=<nil>
}
```

## Performance Benchmark And Alternatives

Marshmallow performs best when dealing with mixed data - when some fields are known and some are unknown.
More info [below](#when-should-i-use-marshmallow).
Other solutions are available for this kind of use case, each solution is explained and documented in the link below.
The full benchmark test can be found
[here](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go).

|Benchmark|(1)|(2)|(3)|(4)|
|--|--|--|--|--|
|[unmarshall twice](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L40)|228693|5164 ns/op|1640 B/op|51 allocs/op|
|[raw map](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L66)|232236|5116 ns/op|2296 B/op|53 allocs/op|
|[go codec](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L121)|388442|3077 ns/op|2512 B/op|37 allocs/op|
|[marshmallow](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L16)|626168|1853 ns/op|608 B/op|18 allocs/op|
|[marshmallow without populating struct](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L162)|678616|1751 ns/op|608 B/op|18 allocs/op|

![marshmallow performance comparison](https://raw.githubusercontent.com/PerimeterX/marshmallow/e45088ca20d4ea5be4143d418d12da63a68d6dfd/performance-chart.svg)

**Marshmallow provides the best performance (up to X3 faster) while not requiring any extra coding.**
In fact, marshmallow performs as fast as normal `json.Unmarshal` call, however, such a call causes loss of data for all
the fields that did not match the given struct. With marshmallow you never loose any data.

|Benchmark|(1)|(2)|(3)|(4)|
|--|--|--|--|--|
|[marshmallow](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L16)|626168|1853 ns/op|608 B/op|18 allocs/op|
|[native library](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L143)|652106|1845 ns/op|304 B/op|11 allocs/op|
|[marshmallow without populating struct](https://github.com/PerimeterX/marshmallow/blob/8c5bba9e6dc0033f4324eca554737089a99f6e5e/benchmark_test.go#L162)|678616|1751 ns/op|608 B/op|18 allocs/op|

## When Should I Use Marshmallow

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

## API

Marshmallow exposes two main API functions - 
[Unmarshal](https://github.com/PerimeterX/marshmallow/blob/0e0218ab860be8a4b5f57f5ff239f281c250c5da/unmarshal.go#L27)
and
[UnmarshalFromJSONMap](https://github.com/PerimeterX/marshmallow/blob/0e0218ab860be8a4b5f57f5ff239f281c250c5da/unmarshal_from_json_map.go#L37).
Each of them can operate in three possible [modes](https://github.com/PerimeterX/marshmallow/blob/0e0218ab860be8a4b5f57f5ff239f281c250c5da/options.go#L30),
and allow setting [skipPopulateStruct](https://github.com/PerimeterX/marshmallow/blob/0e0218ab860be8a4b5f57f5ff239f281c250c5da/options.go#L41) mode.

Marshmallow also supports caching of refection information using 
[EnableCache](https://github.com/PerimeterX/marshmallow/blob/d3500aa5b0f330942b178b155da933c035dd3906/cache.go#L40)
and
[EnableCustomCache](https://github.com/PerimeterX/marshmallow/blob/d3500aa5b0f330942b178b155da933c035dd3906/cache.go#L35).
