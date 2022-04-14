// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-test/deep"
	"github.com/mailru/easyjson/jlexer"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestUnmarshalInputVariations(t *testing.T) {
	EnableCache(&sync.Map{})
	tests := []struct {
		name                string
		mode                Mode
		expectedErr         bool
		expectedResult      bool
		structModifier      func(*parentStruct)
		inputMapModifier    func(map[string]interface{})
		expectedMapModifier func(map[string]interface{})
	}{
		{
			name:                "ModeFailOnFirstError_happy_flow",
			mode:                ModeFailOnFirstError,
			expectedErr:         false,
			expectedResult:      true,
			structModifier:      nil,
			inputMapModifier:    nil,
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_zero_struct_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier:    nil,
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_null_on_struct",
			mode:           ModeFailOnFirstError,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = nil
			},
		},
		{
			name:           "ModeFailOnFirstError_null_on_string",
			mode:           ModeFailOnFirstError,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field1"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField1 = ""
			},
		},
		{
			name:           "ModeFailOnFirstError_null_on_slice",
			mode:           ModeFailOnFirstError,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field30"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField30 = nil
			},
		},
		{
			name:           "ModeFailOnFirstError_null_on_array",
			mode:           ModeFailOnFirstError,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field31"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField31 = [4]string{}
			},
		},
		{
			name:           "ModeFailOnFirstError_null_on_map",
			mode:           ModeFailOnFirstError,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: nil,
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = nil
			},
		},
		{
			name:           "ModeFailOnFirstError_invalid_struct_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_struct_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field2"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_slice_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_array_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_ptr_slice_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field5"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_ptr_array_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field6"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_primitive_map_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_struct_map_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field8"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_struct_ptr_map_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field9"] = 12
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_string_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field1"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_bool_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field2"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field3"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int8_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field4"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int16_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field5"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int32_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field6"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int64_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field7"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field8"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint8_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field9"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint16_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field10"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint32_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field11"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint64_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field12"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_float32_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field13"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_float64_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field14"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_string_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field15"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_bool_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field16"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field17"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int8_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field18"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int16_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field19"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int32_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field20"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_int64_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field21"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field22"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint8_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field23"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint16_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field24"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint32_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field25"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_uint64_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field26"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_float32_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field27"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_float64_ptr_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field28"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_string_slice_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field30"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_string_array_value",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field31"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_slice_element",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = []interface{}{nil, "foo", nil, nil}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_array_element",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = []interface{}{nil, "foo", nil, nil}
			},
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOnFirstError_invalid_map_entry",
			mode:           ModeFailOnFirstError,
			expectedErr:    true,
			expectedResult: false,
			structModifier: nil,
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = map[string]interface{}{"foo": "a", "goo": 12, "boo": "c"}
			},
			expectedMapModifier: nil,
		},
		{
			name:                "ModeAllowMultipleErrors_happy_flow",
			mode:                ModeAllowMultipleErrors,
			expectedErr:         false,
			expectedResult:      true,
			structModifier:      nil,
			inputMapModifier:    nil,
			expectedMapModifier: nil,
		},
		{
			name:           "ModeAllowMultipleErrors_zero_struct_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier:    nil,
			expectedMapModifier: nil,
		},
		{
			name:           "ModeAllowMultipleErrors_null_on_struct",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_null_on_string",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field1"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField1 = ""
			},
		},
		{
			name:           "ModeAllowMultipleErrors_null_on_slice",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field30"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField30 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_null_on_array",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field31"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField31 = [4]string{}
			},
		},
		{
			name:           "ModeAllowMultipleErrors_null_on_map",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: nil,
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_struct_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field1")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_struct_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField2 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field2"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field2")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_slice_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField3 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field3")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_array_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField4 = [4]childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field4")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_ptr_slice_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField5 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field5"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field5")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_ptr_array_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField6 = [4]*childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field6"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field6")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_primitive_map_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field7")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_struct_map_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField8 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field8"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field8")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_struct_ptr_map_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField9 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field9"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				delete(m, "parent_field9")
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_string_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField1 = ""
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field1"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField1 = ""
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_bool_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField2 = false
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field2"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField2 = false
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField3 = 0
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field3"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField3 = 0
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int8_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField4 = int8(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field4"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField4 = int8(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int16_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField5 = int16(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field5"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField5 = int16(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int32_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField6 = int32(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field6"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField6 = int32(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int64_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField7 = int64(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field7"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField7 = int64(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField8 = uint(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field8"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField8 = uint(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint8_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField9 = uint8(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field9"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField9 = uint8(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint16_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField10 = uint16(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field10"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField10 = uint16(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint32_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField11 = uint32(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field11"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField11 = uint32(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint64_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField12 = uint64(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field12"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField12 = uint64(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_float32_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField13 = float32(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field13"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField13 = float32(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_float64_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField14 = float64(0)
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field14"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField14 = float64(0)
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_string_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField15 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field15"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField15 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_bool_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField16 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field16"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField16 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField17 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field17"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField17 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int8_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField18 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field18"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField18 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int16_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField19 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field19"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField19 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int32_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField20 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field20"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField20 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_int64_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField21 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field21"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField21 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField22 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field22"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField22 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint8_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField23 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field23"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField23 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint16_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField24 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field24"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField24 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint32_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField25 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field25"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField25 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_uint64_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField26 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field26"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField26 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_float32_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField27 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field27"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField27 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_float64_ptr_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField28 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field28"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField28 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_string_slice_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField30 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field30"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField30 = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_string_array_value",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1.ChildField31 = [4]string{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field31"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField31 = [4]string{}
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_slice_element",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField3 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = []interface{}{nil, "foo", nil, nil}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_array_element",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField4 = [4]childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = []interface{}{nil, "foo", nil, nil}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = nil
			},
		},
		{
			name:           "ModeAllowMultipleErrors_invalid_map_entry",
			mode:           ModeAllowMultipleErrors,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = map[string]interface{}{"foo": "a", "goo": 12, "boo": "c"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = nil
			},
		},
		{
			name:                "ModeFailOverToOriginalValue_happy_flow",
			mode:                ModeFailOverToOriginalValue,
			expectedErr:         false,
			expectedResult:      true,
			structModifier:      nil,
			inputMapModifier:    nil,
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOverToOriginalValue_zero_struct_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier:    nil,
			expectedMapModifier: nil,
		},
		{
			name:           "ModeFailOverToOriginalValue_null_on_struct",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = nil
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_null_on_string",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field1"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField1 = ""
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_null_on_slice",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field30"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField30 = nil
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_null_on_array",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field31"] = nil
			},
			expectedMapModifier: func(m map[string]interface{}) {
				c := m["parent_field1"].(childStruct)
				c.ChildField31 = [4]string{}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_null_on_map",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    false,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: nil,
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = nil
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_struct_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_struct_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField2 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field2"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field2"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_slice_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField3 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_array_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField4 = [4]childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_ptr_slice_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField5 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field5"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field5"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_ptr_array_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField6 = [4]*childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field6"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field6"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_primitive_map_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_struct_map_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField8 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field8"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field8"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_struct_ptr_map_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField9 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field9"] = 12
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field9"] = float64(12)
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_string_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field1"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field1"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_bool_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field2"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field2"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field3"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field3"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int8_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field4"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field4"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int16_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field5"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field5"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int32_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field6"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field6"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int64_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field7"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field7"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field8"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field8"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint8_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field9"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field9"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint16_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field10"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field10"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint32_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field11"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field11"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint64_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field12"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field12"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_float32_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field13"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field13"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_float64_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field14"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field14"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_string_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field15"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field15"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_bool_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field16"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field16"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field17"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field17"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int8_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field18"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field18"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int16_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field19"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field19"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int32_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field20"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field20"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_int64_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field21"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field21"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field22"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field22"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint8_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field23"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field23"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint16_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field24"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field24"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint32_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field25"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field25"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_uint64_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field26"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field26"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_float32_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field27"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field27"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_float64_ptr_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field28"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field28"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_string_slice_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field30"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field30"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_string_array_value",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField1 = childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field1"].(map[string]interface{})["child_field31"] = map[string]interface{}{"foo": "boo"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field1"] = toMap(m["parent_field1"])
				m["parent_field1"].(map[string]interface{})["child_field31"] = map[string]interface{}{"foo": "boo"}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_slice_element",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField3 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = []interface{}{nil, "foo", nil, nil}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field3"] = []interface{}{childStruct{}, "foo", nil, nil}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_array_element",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField4 = [4]childStruct{}
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = []interface{}{nil, "foo", nil, nil}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field4"] = []interface{}{childStruct{}, "foo", nil, nil}
			},
		},
		{
			name:           "ModeFailOverToOriginalValue_invalid_map_entry",
			mode:           ModeFailOverToOriginalValue,
			expectedErr:    true,
			expectedResult: true,
			structModifier: func(p *parentStruct) {
				p.ParentField7 = nil
			},
			inputMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = map[string]interface{}{"foo": "a", "goo": 12, "boo": "c"}
			},
			expectedMapModifier: func(m map[string]interface{}) {
				m["parent_field7"] = map[string]interface{}{"foo": "a", "goo": float64(12), "boo": "c"}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedStruct := buildParentStruct()
			if tt.structModifier != nil {
				tt.structModifier(expectedStruct)
			}
			input := toMap(expectedStruct)
			for k, v := range extraData {
				input[k] = v
			}
			if tt.inputMapModifier != nil {
				tt.inputMapModifier(input)
			}
			data, err := json.Marshal(input)
			if err != nil {
				t.Errorf("could not marshal parent: %v", err)
				return
			}
			actualStruct := &parentStruct{}
			actualMap, err := Unmarshal(data, actualStruct, WithMode(tt.mode))
			if (err != nil) != tt.expectedErr {
				t.Errorf("Unmarshal() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if tt.expectedResult {
				expectedStruct.ParentField10.CustomField = "UnmarshalJSON called"
				expectedStruct.ParentField11.CustomField = "UnmarshalJSON called"
				if diff := deep.Equal(actualStruct, expectedStruct); diff != nil {
					t.Errorf("Unmarshal() struct mismatch (actual, expected):\n%s", strings.Join(diff, "\n"))
				}
				expectedMap := make(map[string]interface{})
				for k, v := range extraData {
					expectedMap[k] = v
				}
				structValue := reflectStructValue(actualStruct)
				for name, fieldIdx := range mapStructFields(actualStruct) {
					field := structValue.Field(fieldIdx)
					expectedMap[name] = field.Interface()
				}
				if tt.expectedMapModifier != nil {
					tt.expectedMapModifier(expectedMap)
				}
				if tt.mode == ModeFailOverToOriginalValue {
					normalizeMapTypes(actualMap)
				}
				if diff := deep.Equal(actualMap, expectedMap); diff != nil {
					t.Errorf("Unmarshal() map mismatch (actual, expected):\n%s", strings.Join(diff, "\n"))
				}
			} else {
				if reflect.DeepEqual(actualStruct, expectedStruct) {
					t.Error("Unmarshal() expected parsing to break before finished")
				}
				if actualMap != nil {
					t.Errorf("Unmarshal() expected actual map to not exist")
				}
			}
		})
	}
}

func TestUnmarshalSpecialInput(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		v            interface{}
		mode         Mode
		result       bool
		errValidator func(error) bool
	}{
		{
			name:   "invalid_input",
			data:   []byte(`12`),
			v:      &parentStruct{},
			mode:   ModeFailOnFirstError,
			result: false,
			errValidator: func(err error) bool {
				return err == ErrInvalidInput
			},
		},
		{
			name:   "invalid_value",
			data:   []byte(`{"field":""}`),
			v:      "",
			mode:   ModeFailOnFirstError,
			result: false,
			errValidator: func(err error) bool {
				return err == ErrInvalidValue
			},
		},
		{
			name:   "null_input",
			data:   []byte(`null`),
			v:      &parentStruct{},
			mode:   ModeFailOnFirstError,
			result: true,
			errValidator: func(err error) bool {
				return err == nil
			},
		},
		{
			name:   "ModeFailOnFirstError_custom_unmarshal_failing",
			data:   []byte(`{"field":""}`),
			v:      &failingCustomUnmarshalerParent{},
			mode:   ModeFailOnFirstError,
			result: false,
			errValidator: func(err error) bool {
				e, ok := err.(*jlexer.LexerError)
				if !ok {
					return false
				}
				return e.Reason == "failing"
			},
		},
		{
			name:   "ModeAllowMultipleErrors_custom_unmarshal_failing",
			data:   []byte(`{"field":""}`),
			v:      &failingCustomUnmarshalerParent{},
			mode:   ModeAllowMultipleErrors,
			result: true,
			errValidator: func(err error) bool {
				e, ok := err.(*MultipleLexerError)
				if !ok {
					return false
				}
				if len(e.Errors) != 1 {
					return false
				}
				return e.Errors[0].Reason == "failing"
			},
		},
		{
			name:   "ModeFailOverToOriginalValue_custom_unmarshal_failing",
			data:   []byte(`{"field":""}`),
			v:      &failingCustomUnmarshalerParent{},
			mode:   ModeFailOverToOriginalValue,
			result: true,
			errValidator: func(err error) bool {
				e, ok := err.(*MultipleLexerError)
				if !ok {
					return false
				}
				if len(e.Errors) != 1 {
					return false
				}
				return e.Errors[0].Reason == "failing"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.data, tt.v, WithMode(tt.mode))
			if !tt.errValidator(err) {
				t.Errorf("Unmarshal() unexpected error = %v", err)
				return
			}
			if tt.result {
				if got == nil {
					t.Error("Unmarshal() expected result exists")
					return
				}
			} else {
				if got != nil {
					t.Error("Unmarshal() expected result not exists")
					return
				}
			}
		})
	}
}

var extraData = map[string]interface{}{
	"extra1": "foo",
	"extra2": float64(12),
	"extra3": true,
	"extra4": []interface{}{"1", false},
}

func buildParentStruct() *parentStruct {
	return &parentStruct{
		ParentField1:  *buildChildStruct(),
		ParentField2:  buildChildStruct(),
		ParentField3:  []childStruct{*buildChildStruct()},
		ParentField4:  [4]childStruct{*buildChildStruct()},
		ParentField5:  []*childStruct{buildChildStruct(), nil},
		ParentField6:  [4]*childStruct{buildChildStruct()},
		ParentField7:  map[string]string{"f6-key-1": "f6-value-1", "f6-key-2": "f6-value-2"},
		ParentField8:  map[string]childStruct{"f7-key-1": *buildChildStruct()},
		ParentField9:  map[string]*childStruct{"f8-key-1": buildChildStruct()},
		ParentField10: customUnmarshaler{CustomField: "ignore this"},
		ParentField11: &customUnmarshaler{CustomField: "ignore this too"},
	}
}

func buildChildStruct() *childStruct {
	f15 := "field15"
	f16 := true
	f17 := 17
	f18 := int8(18)
	f19 := int16(19)
	f20 := int32(20)
	f21 := int64(21)
	f22 := uint(22)
	f23 := uint8(23)
	f24 := uint16(24)
	f25 := uint32(25)
	f26 := uint64(26)
	f27 := float32(27.7)
	f28 := 28.8
	return &childStruct{
		ChildField1:  "field1",
		ChildField2:  true,
		ChildField3:  3,
		ChildField4:  4,
		ChildField5:  5,
		ChildField6:  6,
		ChildField7:  7,
		ChildField8:  8,
		ChildField9:  9,
		ChildField10: 10,
		ChildField11: 11,
		ChildField12: 12,
		ChildField13: 13.3,
		ChildField14: 14.4,
		ChildField15: &f15,
		ChildField16: &f16,
		ChildField17: &f17,
		ChildField18: &f18,
		ChildField19: &f19,
		ChildField20: &f20,
		ChildField21: &f21,
		ChildField22: &f22,
		ChildField23: &f23,
		ChildField24: &f24,
		ChildField25: &f25,
		ChildField26: &f26,
		ChildField27: &f27,
		ChildField28: &f28,
		ChildField29: "interface",
		ChildField30: []string{"f30-1", "f30-2"},
		ChildField31: [4]string{"f31-1", "f31-2", "f31-3", "f31-4"},
	}
}

func toMap(value interface{}) map[string]interface{} {
	data, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("could not marshal value to map %v", err))
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal value to map %v", err))
	}
	return result
}

type parentStruct struct {
	ParentField1  childStruct             `json:"parent_field1"`
	ParentField2  *childStruct            `json:"parent_field2"`
	ParentField3  []childStruct           `json:"parent_field3"`
	ParentField4  [4]childStruct          `json:"parent_field4"`
	ParentField5  []*childStruct          `json:"parent_field5"`
	ParentField6  [4]*childStruct         `json:"parent_field6"`
	ParentField7  map[string]string       `json:"parent_field7"`
	ParentField8  map[string]childStruct  `json:"parent_field8"`
	ParentField9  map[string]*childStruct `json:"parent_field9"`
	ParentField10 customUnmarshaler       `json:"parent_field10"`
	ParentField11 *customUnmarshaler      `json:"parent_field11"`
}

type childStruct struct {
	ChildField1  string      `json:"child_field1,omitempty"`
	ChildField2  bool        `json:"child_field2"`
	ChildField3  int         `json:"child_field3"`
	ChildField4  int8        `json:"child_field4"`
	ChildField5  int16       `json:"child_field5"`
	ChildField6  int32       `json:"child_field6"`
	ChildField7  int64       `json:"child_field7"`
	ChildField8  uint        `json:"child_field8"`
	ChildField9  uint8       `json:"child_field9"`
	ChildField10 uint16      `json:"child_field10"`
	ChildField11 uint32      `json:"child_field11"`
	ChildField12 uint64      `json:"child_field12"`
	ChildField13 float32     `json:"child_field13"`
	ChildField14 float64     `json:"child_field14"`
	ChildField15 *string     `json:"child_field15"`
	ChildField16 *bool       `json:"child_field16"`
	ChildField17 *int        `json:"child_field17"`
	ChildField18 *int8       `json:"child_field18"`
	ChildField19 *int16      `json:"child_field19"`
	ChildField20 *int32      `json:"child_field20"`
	ChildField21 *int64      `json:"child_field21"`
	ChildField22 *uint       `json:"child_field22"`
	ChildField23 *uint8      `json:"child_field23"`
	ChildField24 *uint16     `json:"child_field24"`
	ChildField25 *uint32     `json:"child_field25"`
	ChildField26 *uint64     `json:"child_field26"`
	ChildField27 *float32    `json:"child_field27"`
	ChildField28 *float64    `json:"child_field28"`
	ChildField29 interface{} `json:"child_field29"`
	ChildField30 []string    `json:"child_field30"`
	ChildField31 [4]string   `json:"child_field31"`
}

type customUnmarshaler struct {
	CustomField string
}

func (c *customUnmarshaler) UnmarshalJSON([]byte) error {
	*c = customUnmarshaler{CustomField: "UnmarshalJSON called"}
	return nil
}

func (c *customUnmarshaler) UnmarshalJSONFromMap(interface{}) error {
	*c = customUnmarshaler{CustomField: "UnmarshalJSON called"}
	return nil
}

type failingCustomUnmarshalerParent struct {
	Field *failingCustomUnmarshaler `json:"field"`
}

type failingCustomUnmarshaler struct{}

func (c *failingCustomUnmarshaler) UnmarshalJSON([]byte) error {
	return errors.New("failing")
}

func (c *failingCustomUnmarshaler) UnmarshalJSONFromMap(interface{}) error {
	return errors.New("failing")
}

func normalizeMapTypes(m map[string]interface{}) {
	for k, v := range m {
		switch actual := v.(type) {
		case uint:
			m[k] = float64(actual)
		case uint8:
			m[k] = float64(actual)
		case uint16:
			m[k] = float64(actual)
		case uint32:
			m[k] = float64(actual)
		case uint64:
			m[k] = float64(actual)
		case int:
			m[k] = float64(actual)
		case int8:
			m[k] = float64(actual)
		case int16:
			m[k] = float64(actual)
		case int32:
			m[k] = float64(actual)
		case int64:
			m[k] = float64(actual)
		case float32:
			m[k] = float64(actual)
		case []string:
			data := make([]interface{}, len(actual))
			for i, item := range actual {
				data[i] = item
			}
			m[k] = data
		case [4]string:
			data := make([]interface{}, len(actual))
			for i, item := range actual {
				data[i] = item
			}
			m[k] = data
		case map[string]interface{}:
			normalizeMapTypes(actual)
		}
	}
}
