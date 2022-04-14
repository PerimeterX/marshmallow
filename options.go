// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

func WithMode(mode Mode) UnmarshalOption {
	return func(options *unmarshalOptions) {
		options.mode = mode
	}
}

type UnmarshalOption func(*unmarshalOptions)

type unmarshalOptions struct {
	mode Mode
}

type Mode uint8

const (
	ModeFailOnFirstError Mode = iota
	ModeAllowMultipleErrors
	ModeFailOverToOriginalValue
)

func buildUnmarshalOptions(options []UnmarshalOption) *unmarshalOptions {
	result := &unmarshalOptions{}
	for _, option := range options {
		option(result)
	}
	return result
}
