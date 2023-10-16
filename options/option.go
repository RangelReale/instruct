package options

import "github.com/rrgmc/instruct"

// AnyOption states that the option can be used in any instruct.Decoder function.
type AnyOption[IT any, DC instruct.DecodeContext] interface {
	isAnyOption()
}

// AnyTypeOption states that the option can be used in any instruct.TypeDecoder function.
type AnyTypeOption[IT any, DC instruct.DecodeContext] interface {
	isAnyTypeOption()
}

// DefaultOption states that the option can be used in the instruct.NewDecoder function.
type DefaultOption[IT any, DC instruct.DecodeContext, DOPTT any] interface {
	AnyOption[IT, DC]
	ApplyDefaultOption(*DOPTT)
}

// TypeDefaultOption states that the option can be used in the instruct.NewTypeDecoder function.
type TypeDefaultOption[IT any, DC instruct.DecodeContext, TOPTT any] interface {
	AnyTypeOption[IT, DC]
	ApplyTypeDefaultOption(*TOPTT)
}

// DecodeOption states that the option can be used in the [instruct.Decoder.Decode] function.
type DecodeOption[IT any, DC instruct.DecodeContext, COPTT any] interface {
	AnyOption[IT, DC]
	ApplyDecodeOption(*COPTT)
}

// TypeDecodeOption states that the option can be used in the [instruct.TypeDecoder.Decode] function.
type TypeDecodeOption[IT any, DC instruct.DecodeContext, COPTT any] interface {
	AnyTypeOption[IT, DC]
	ApplyTypeDecodeOption(*COPTT)
}

// DefaultAndTypeDefaultOption states that the option can be used in the instruct.NewDecoder and
// instruct.NewTypeDecoder functions.
type DefaultAndTypeDefaultOption[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any] interface {
	DefaultOption[IT, DC, DOPTT]
	TypeDefaultOption[IT, DC, TOPTT]
}

// DefaultAndDecodeOption states that the option can be used in the instruct.NewDecoder and
// [instruct.Decoder.Decode] functions.
type DefaultAndDecodeOption[IT any, DC instruct.DecodeContext, DOPTT any, COPTT any] interface {
	DefaultOption[IT, DC, DOPTT]
	DecodeOption[IT, DC, COPTT]
}

// TypeDefaultAndTypeDecodeOption states that the option can be used in the instruct.NewTypeDecoder and
// [instruct.TypeDecoder.Decode] functions.
type TypeDefaultAndTypeDecodeOption[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any] interface {
	TypeDefaultOption[IT, DC, TOPTT]
	TypeDecodeOption[IT, DC, COPTT]
}

// DefaultAndTypeDefaultDecodeOption states that the option can be used in the [instruct.Decoder.Decode] and
// [instruct.TypeDecoder.Decode] functions.
type DefaultAndTypeDefaultDecodeOption[IT any, DC instruct.DecodeContext, DCOPTT any, TCOPTT any] interface {
	DecodeOption[IT, DC, DCOPTT]
	TypeDecodeOption[IT, DC, TCOPTT]
}

// TypeDefaultAndDecodeOption states that the option can be used in the instruct.NewTypeDecoder and
// [instruct.Decoder.Decode] functions.
type TypeDefaultAndDecodeOption[IT any, DC instruct.DecodeContext, TOPTT any, DCOPTT any] interface {
	TypeDefaultOption[IT, DC, TOPTT]
	DecodeOption[IT, DC, DCOPTT]
}

// FullOption states that the option can be used in all New functions and Decode methods.
type FullOption[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any, DCOPTT any, TCOPTT any] interface {
	DefaultOption[IT, DC, DOPTT]
	TypeDefaultOption[IT, DC, TOPTT]
	DecodeOption[IT, DC, DCOPTT]
	TypeDecodeOption[IT, DC, TCOPTT]
}
