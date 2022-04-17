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

#### Where does marshmallow shine
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

## Alternatives
Other solution available for this kind of use case:
- Unmarshal twice: once into a struct and a second time into a map. This is fully
native and requires no external dependencies. However, it obviously has huge implications
over performance. This approach will be useful in case performance does not matter,
and you do not wish to import any external dependencies.
- Unmarshal into a raw map and then - [example](https://stackoverflow.com/a/33499066/1932186).
This method will be useful if you are willing to write more code to boost performance
just by a bit and still avoid using external dependencies.
- Use `go/codec` or other libraries - [example](https://stackoverflow.com/a/33499861/1932186).
This will boost a bit more performance but require explicit coding to hook struct fields 
into the map.

## Performance


## API
API, options and cache

