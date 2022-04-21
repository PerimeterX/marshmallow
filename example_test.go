package marshmallow_test

import (
	"fmt"
	"github.com/perimeterx/marshmallow"
)

type exampleStruct struct {
	Foo string `json:"foo"`
	Boo []int  `json:"boo"`
}

func ExampleUnmarshal() {
	// enable marshmallow cache to boost up performance by reusing field type information.
	marshmallow.EnableCache()

	// unmarshal with mode marshmallow.ModeFailOnFirstError and valid value
	// this will finish unmarshalling and return a nil err
	v := exampleStruct{}
	result, err := marshmallow.Unmarshal([]byte(`{"foo":"bar","boo":[1,2,3]}`), &v)
	fmt.Printf("ModeFailOnFirstError and valid value: v=%+v, result=%+v, err=%T\n", v, result, err)

	// unmarshal with mode marshmallow.ModeFailOnFirstError and invalid value
	// this will return nil result and an error
	v = exampleStruct{}
	result, err = marshmallow.Unmarshal([]byte(`{"foo":2,"boo":[1,2,3]}`), &v)
	fmt.Printf("ModeFailOnFirstError and invalid value: result=%+v, err=%T\n", result, err)

	// unmarshal with mode marshmallow.ModeAllowMultipleErrors and valid value
	// this will finish unmarshalling and return a nil err
	v = exampleStruct{}
	result, err = marshmallow.Unmarshal([]byte(`{"foo":"bar","boo":[1,2,3]}`), &v,
		marshmallow.WithMode(marshmallow.ModeAllowMultipleErrors))
	fmt.Printf("ModeAllowMultipleErrors and valid value: v=%+v, result=%+v, err=%T\n", v, result, err)

	// unmarshal with mode marshmallow.ModeAllowMultipleErrors and invalid value
	// this will return a partially populated result and an error
	v = exampleStruct{}
	result, err = marshmallow.Unmarshal([]byte(`{"foo":2,"boo":[1,2,3]}`), &v,
		marshmallow.WithMode(marshmallow.ModeAllowMultipleErrors))
	fmt.Printf("ModeAllowMultipleErrors and invalid value: result=%+v, err=%T\n", result, err)

	// unmarshal with mode marshmallow.ModeFailOverToOriginalValue and valid value
	// this will finish unmarshalling and return a nil err
	v = exampleStruct{}
	result, err = marshmallow.Unmarshal([]byte(`{"foo":"bar","boo":[1,2,3]}`), &v,
		marshmallow.WithMode(marshmallow.ModeFailOverToOriginalValue))
	fmt.Printf("ModeFailOverToOriginalValue and valid value: v=%+v, result=%+v, err=%T\n", v, result, err)

	// unmarshal with mode marshmallow.ModeFailOverToOriginalValue and invalid value
	// this will return a fully unmarshalled result, failing to the original invalid values, and an error
	v = exampleStruct{}
	result, err = marshmallow.Unmarshal([]byte(`{"foo":2,"boo":[1,2,3]}`), &v,
		marshmallow.WithMode(marshmallow.ModeFailOverToOriginalValue))
	fmt.Printf("ModeFailOverToOriginalValue and invalid value: result=%+v, err=%T\n", result, err)
	// Output:
	// ModeFailOnFirstError and valid value: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar], err=<nil>
	// ModeFailOnFirstError and invalid value: result=map[], err=*jlexer.LexerError
	// ModeAllowMultipleErrors and valid value: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar], err=<nil>
	// ModeAllowMultipleErrors and invalid value: result=map[boo:[1 2 3]], err=*marshmallow.MultipleLexerError
	// ModeFailOverToOriginalValue and valid value: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar], err=<nil>
	// ModeFailOverToOriginalValue and invalid value: result=map[boo:[1 2 3] foo:2], err=*marshmallow.MultipleLexerError
}

func ExampleUnmarshalFromJSONMap() {
	// enable marshmallow cache to boost up performance by reusing field type information.
	marshmallow.EnableCache()

	// unmarshal with mode marshmallow.ModeFailOnFirstError and valid value
	// this will finish unmarshalling and return a nil err
	v := exampleStruct{}
	result, err := marshmallow.UnmarshalFromJSONMap(
		map[string]interface{}{"foo": "bar", "boo": []interface{}{float64(1), float64(2), float64(3)}}, &v)
	fmt.Printf("ModeFailOnFirstError and valid value: v=%+v, result=%+v, err=%T\n", v, result, err)

	// unmarshal with mode marshmallow.ModeFailOnFirstError and invalid value
	// this will return nil result and an error
	v = exampleStruct{}
	result, err = marshmallow.UnmarshalFromJSONMap(
		map[string]interface{}{"foo": float64(2), "boo": []interface{}{float64(1), float64(2), float64(3)}}, &v)
	fmt.Printf("ModeFailOnFirstError and invalid value: result=%+v, err=%T\n", result, err)

	// unmarshal with mode marshmallow.ModeAllowMultipleErrors and valid value
	// this will finish unmarshalling and return a nil err
	v = exampleStruct{}
	result, err = marshmallow.UnmarshalFromJSONMap(
		map[string]interface{}{"foo": "bar", "boo": []interface{}{float64(1), float64(2), float64(3)}}, &v,
		marshmallow.WithMode(marshmallow.ModeAllowMultipleErrors))
	fmt.Printf("ModeAllowMultipleErrors and valid value: v=%+v, result=%+v, err=%T\n", v, result, err)

	// unmarshal with mode marshmallow.ModeAllowMultipleErrors and invalid value
	// this will return a partially populated result and an error
	v = exampleStruct{}
	result, err = marshmallow.UnmarshalFromJSONMap(
		map[string]interface{}{"foo": float64(2), "boo": []interface{}{float64(1), float64(2), float64(3)}}, &v,
		marshmallow.WithMode(marshmallow.ModeAllowMultipleErrors))
	fmt.Printf("ModeAllowMultipleErrors and invalid value: result=%+v, err=%T\n", result, err)

	// unmarshal with mode marshmallow.ModeFailOverToOriginalValue and valid value
	// this will finish unmarshalling and return a nil err
	v = exampleStruct{}
	result, err = marshmallow.UnmarshalFromJSONMap(
		map[string]interface{}{"foo": "bar", "boo": []interface{}{float64(1), float64(2), float64(3)}}, &v,
		marshmallow.WithMode(marshmallow.ModeFailOverToOriginalValue))
	fmt.Printf("ModeFailOverToOriginalValue and valid value: v=%+v, result=%+v, err=%T\n", v, result, err)

	// unmarshal with mode marshmallow.ModeFailOverToOriginalValue and invalid value
	// this will return a fully unmarshalled result, failing to the original invalid values, and an error
	v = exampleStruct{}
	result, err = marshmallow.UnmarshalFromJSONMap(
		map[string]interface{}{"foo": float64(2), "boo": []interface{}{float64(1), float64(2), float64(3)}}, &v,
		marshmallow.WithMode(marshmallow.ModeFailOverToOriginalValue))
	fmt.Printf("ModeFailOverToOriginalValue and invalid value: result=%+v, err=%T\n", result, err)
	// Output:
	// ModeFailOnFirstError and valid value: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar], err=<nil>
	// ModeFailOnFirstError and invalid value: result=map[], err=*marshmallow.ParseError
	// ModeAllowMultipleErrors and valid value: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar], err=<nil>
	// ModeAllowMultipleErrors and invalid value: result=map[boo:[1 2 3]], err=*marshmallow.MultipleError
	// ModeFailOverToOriginalValue and valid value: v={Foo:bar Boo:[1 2 3]}, result=map[boo:[1 2 3] foo:bar], err=<nil>
	// ModeFailOverToOriginalValue and invalid value: result=map[boo:[1 2 3] foo:2], err=*marshmallow.MultipleError
}
