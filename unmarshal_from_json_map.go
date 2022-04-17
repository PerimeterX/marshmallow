// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

import (
	"reflect"
)

// UnmarshalerFromJSONMap is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// UnmarshalerFromJSONMap should use the value in data to perform the
// unmarshalling from map.
type UnmarshalerFromJSONMap interface {
	UnmarshalJSONFromMap(data interface{}) error
}

// UnmarshalFromJSONMap parses the JSON map data and stores the values
// in the struct pointed to by v and in the returned map.
// If v is nil or not a pointer to a struct, UnmarshalFromJSONMap returns an ErrInvalidValue.
//
// UnmarshalFromJSONMap follows the rules of json.Unmarshal with the following exceptions:
// - All input fields are stored in the resulting map, including fields that do not exist in the
// struct pointed by v. This allows both to retain all the original input data fields and to access
// them via a map key lookup.
// - UnmarshalFromJSONMap receive a JSON map instead of raw bytes. The given input map is assumed
// to be a JSON map, meaning it should only contain the following types: bool, string, float64,
// []interface, and map[string]interface{}. Other types will cause decoding to return unexpected results.
// - UnmarshalFromJSONMap only operates on struct values. It will reject all other types of v by
// returning ErrInvalidValue.
// - UnmarshalFromJSONMap supports three types of Mode values. Each mode is self documented and affects
// how UnmarshalFromJSONMap behaves.
func UnmarshalFromJSONMap(data map[string]interface{}, v interface{}, options ...UnmarshalOption) (map[string]interface{}, error) {
	if !isValidValue(v) {
		return nil, ErrInvalidValue
	}
	opts := buildUnmarshalOptions(options)
	d := &mapDecoder{options: opts}
	result := make(map[string]interface{})
	if data != nil {
		d.populateStruct(nil, data, v, result)
	}
	if opts.mode == ModeAllowMultipleErrors || opts.mode == ModeFailOverToOriginalValue {
		if len(d.errs) == 0 {
			return result, nil
		}
		return result, &MultipleError{Errors: d.errs}
	}
	if d.err != nil {
		return nil, d.err
	}
	return result, nil
}

var unmarshalerFromJSONMapType = reflect.TypeOf((*UnmarshalerFromJSONMap)(nil)).Elem()

type mapDecoder struct {
	options *unmarshalOptions
	err     error
	errs    []error
}

func (m *mapDecoder) populateStruct(path []string, data map[string]interface{}, structInstance interface{}, result map[string]interface{}) (interface{}, bool) {
	structValue := reflectStructValue(structInstance)
	fields := mapStructFields(structInstance)
	for key, inputValue := range data {
		fieldIdx, exists := fields[key]
		if exists {
			field := structValue.Field(fieldIdx)
			value, isValidType := m.valueByReflectType(append(path, key), inputValue, field.Type(), false)
			if isValidType {
				if !m.options.skipPopulateStruct || result == nil {
					assignValue(field, value)
				}
				if result != nil {
					result[key] = value
				}
			} else {
				switch m.options.mode {
				case ModeFailOnFirstError:
					return nil, false
				case ModeFailOverToOriginalValue:
					if result != nil {
						result[key] = value
					} else {
						return data, false
					}
				}
			}
		} else {
			if result != nil {
				result[key] = inputValue
			}
		}
	}
	return structInstance, true
}

func (m *mapDecoder) valueByReflectType(path []string, v interface{}, t reflect.Type, isPtr bool) (interface{}, bool) {
	if t.Implements(unmarshalerFromJSONMapType) {
		result := reflect.New(t.Elem()).Interface()
		m.valueFromCustomUnmarshaler(v, result.(UnmarshalerFromJSONMap))
		return result, true
	}
	if reflect.PtrTo(t).Implements(unmarshalerFromJSONMapType) {
		value := reflect.New(t)
		m.valueFromCustomUnmarshaler(v, value.Interface().(UnmarshalerFromJSONMap))
		return value.Elem().Interface(), true
	}
	kind := t.Kind()
	if converter := primitiveConverters[kind]; converter != nil {
		if v == nil {
			if isPtr || kind == reflect.Interface {
				return v, true
			}
			return reflect.Zero(t).Interface(), true
		}
		converted, ok := converter(v)
		if !ok {
			m.addError(newUnexpectedTypeParseError(t, path))
			return v, false
		}
		return converted, true
	}
	switch kind {
	case reflect.Slice:
		return m.buildSlice(path, v, t)
	case reflect.Array:
		return m.buildArray(path, v, t)
	case reflect.Map:
		return m.buildMap(path, v, t)
	case reflect.Struct:
		value, valid := m.buildStruct(path, v, t)
		if value == nil {
			return nil, valid
		}
		if !valid {
			return value, false
		}
		return reflect.ValueOf(value).Elem().Interface(), valid
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Struct {
			return m.buildStruct(path, v, t.Elem())
		}
		value, valid := m.valueByReflectType(path, v, t.Elem(), true)
		if value == nil {
			return nil, valid
		}
		if !valid {
			return value, false
		}
		result := reflect.New(reflect.TypeOf(value))
		result.Elem().Set(reflect.ValueOf(value))
		return result.Interface(), valid
	}
	m.addError(newUnsupportedTypeParseError(t, path))
	return nil, false
}

func (m *mapDecoder) buildSlice(path []string, v interface{}, sliceType reflect.Type) (interface{}, bool) {
	if v == nil {
		return nil, true
	}
	arr, ok := v.([]interface{})
	if !ok {
		m.addError(newUnexpectedTypeParseError(sliceType, path))
		return v, false
	}
	elemType := sliceType.Elem()
	var sliceValue reflect.Value
	if len(arr) > 0 {
		sliceValue = reflect.MakeSlice(sliceType, 0, 4)
	} else {
		sliceValue = reflect.MakeSlice(sliceType, 0, 0)
	}
	for _, element := range arr {
		current, valid := m.valueByReflectType(path, element, elemType, false)
		if !valid {
			if m.options.mode != ModeFailOverToOriginalValue {
				return nil, true
			}
			return v, true
		}
		sliceValue = reflect.Append(sliceValue, safeReflectValue(elemType, current))
	}
	return sliceValue.Interface(), true
}

func (m *mapDecoder) buildArray(path []string, v interface{}, arrayType reflect.Type) (interface{}, bool) {
	if v == nil {
		return nil, true
	}
	arr, ok := v.([]interface{})
	if !ok {
		m.addError(newUnexpectedTypeParseError(arrayType, path))
		return v, false
	}
	elemType := arrayType.Elem()
	arrayValue := reflect.New(arrayType).Elem()
	for i, element := range arr {
		current, valid := m.valueByReflectType(path, element, elemType, false)
		if !valid {
			if m.options.mode != ModeFailOverToOriginalValue {
				return nil, true
			}
			return v, true
		}
		if current != nil {
			arrayValue.Index(i).Set(reflect.ValueOf(current))
		}
	}
	return arrayValue.Interface(), true
}

func (m *mapDecoder) buildMap(path []string, v interface{}, mapType reflect.Type) (interface{}, bool) {
	if v == nil {
		return nil, true
	}
	mp, ok := v.(map[string]interface{})
	if !ok {
		m.addError(newUnexpectedTypeParseError(mapType, path))
		return v, false
	}
	keyType := mapType.Key()
	valueType := mapType.Elem()
	mapValue := reflect.MakeMap(mapType)
	for inputKey, inputValue := range mp {
		keyPath := append(path, inputKey)
		key, valid := m.valueByReflectType(keyPath, inputKey, keyType, false)
		if !valid {
			if m.options.mode != ModeFailOverToOriginalValue {
				return nil, true
			}
			return v, true
		}
		value, valid := m.valueByReflectType(keyPath, inputValue, valueType, false)
		if !valid {
			if m.options.mode != ModeFailOverToOriginalValue {
				return nil, true
			}
			return v, true
		}
		mapValue.SetMapIndex(safeReflectValue(keyType, key), safeReflectValue(valueType, value))
	}
	return mapValue.Interface(), true
}

func (m *mapDecoder) buildStruct(path []string, v interface{}, structType reflect.Type) (interface{}, bool) {
	if v == nil {
		return nil, true
	}
	mp, ok := v.(map[string]interface{})
	if !ok {
		m.addError(newUnexpectedTypeParseError(structType, path))
		return v, false
	}
	value := reflect.New(structType).Interface()
	return m.populateStruct(path, mp, value, nil)
}

func (m *mapDecoder) valueFromCustomUnmarshaler(data interface{}, unmarshaler UnmarshalerFromJSONMap) {
	err := unmarshaler.UnmarshalJSONFromMap(data)
	if err != nil {
		m.addError(err)
	}
}

func (m *mapDecoder) addError(err error) {
	if m.options.mode == ModeFailOnFirstError {
		m.err = err
	} else {
		m.errs = append(m.errs, err)
	}
}
