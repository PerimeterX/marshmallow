# Marshmallow

Package marshmallow provides a simple API to perform flexible and performant JSON unmarshalling.
Unlike other packages, marshmallow supports unmarshalling of some known and some unknown fields
with zero performance overhead nor extra coding needed. While unmarshalling,
marshmallow allows fully retaining the original data and access it via a typed struct and a
dynamic map.

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
Marshmallow is best suited for use cases where you are interested in all the input data,
but you have predetermined information only about a subset of it.
For instance, if you plan to reference two specific fields from the data,
then iterate all the data and apply some generic logic. How does it look with the native library:
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

There can be two main reasons to have an interest in all the data.
First is when you eventually plan to write or pipe the input data, and you don't
want to lose any of it. Second is where you plan to perform any kind of dynamic
read of the data - this includes iterating it, reading calculated or configured
field names, and others.

## Alternatives and Performance Benchmark
[Full Benchmark](benchmark_test.go)

Other solutions available for this kind of use case, each solution is explained
and documented in the link below.

|Benchmark|(1)|(2)|(3)|(4)|
|--|--|--|--|--|
|[unmarshall twice](https://github.com/PerimeterX/marshmallow/blob/d165df95a46f197a3db895a542333ae971d9a330/benchmark_test.go#L33)|228693|5164 ns/op|1640 B/op|51 allocs/op|
|[raw map](https://github.com/PerimeterX/marshmallow/blob/d165df95a46f197a3db895a542333ae971d9a330/benchmark_test.go#L33)|232236|5116 ns/op|2296 B/op|53 allocs/op|
|[go codec](https://github.com/PerimeterX/marshmallow/blob/d165df95a46f197a3db895a542333ae971d9a330/benchmark_test.go#L33)|388442|3077 ns/op|2512 B/op|37 allocs/op|
|[marshmallow](https://github.com/PerimeterX/marshmallow/blob/d165df95a46f197a3db895a542333ae971d9a330/benchmark_test.go#L33)|626168|1853 ns/op|608 B/op|18 allocs/op|
|[marshmallow without populating struct](https://github.com/PerimeterX/marshmallow/blob/d165df95a46f197a3db895a542333ae971d9a330/benchmark_test.go#L33)|678616|1751 ns/op|608 B/op|18 allocs/op|

**Marshmallow provides the best performance (up to X3 faster) while not requiring any extra coding.** 

## API
API, options and cache

