package options

import (
	"github.com/RangelReale/instruct"
)

// DefaultOption + TypeDefaultOption

type DefaultOptionImpl[IT any, DC instruct.DecodeContext, DOPTT any] struct {
	f func(o *DOPTT)
}

var _ DefaultOption[any, instruct.DecodeContext, any] = (*DefaultOptionImpl[any, instruct.DecodeContext, any])(nil)

func (f DefaultOptionImpl[IT, DC, DOPTT]) isAnyOption() {}

func (f DefaultOptionImpl[IT, DC, DOPTT]) ApplyDefaultOption(o *DOPTT) {
	f.f(o)
}

// DefaultOptionFunc is the option constructor for DefaultOption.
func DefaultOptionFunc[IT any, DC instruct.DecodeContext, DOPTT any](f func(o *DOPTT)) *DefaultOptionImpl[IT, DC, DOPTT] {
	return &DefaultOptionImpl[IT, DC, DOPTT]{f}
}

// TypeDefaultOption

type TypeDefaultOptionImpl[IT any, DC instruct.DecodeContext, TOPTT any] struct {
	tf func(o *TOPTT)
}

var _ TypeDefaultOption[any, instruct.DecodeContext, any] = (*TypeDefaultOptionImpl[any, instruct.DecodeContext, any])(nil)

func (f TypeDefaultOptionImpl[IT, DC, TOPTT]) isAnyTypeOption() {}

func (f TypeDefaultOptionImpl[IT, DC, TOPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

// TypeDefaultOptionFunc is the option constructor for TypeDefaultOption.
func TypeDefaultOptionFunc[IT any, DC instruct.DecodeContext, TOPTT any](tf func(o *TOPTT)) *TypeDefaultOptionImpl[IT, DC, TOPTT] {
	return &TypeDefaultOptionImpl[IT, DC, TOPTT]{tf}
}

// DecodeOption

type DecodeOptionImpl[IT any, DC instruct.DecodeContext, COPTT any] struct {
	f func(o *COPTT)
}

var _ DecodeOption[any, instruct.DecodeContext, any] = (*DecodeOptionImpl[any, instruct.DecodeContext, any])(nil)

func (f DecodeOptionImpl[IT, DC, COPTT]) isAnyOption() {}

func (f DecodeOptionImpl[IT, DC, COPTT]) ApplyDecodeOption(o *COPTT) {
	f.f(o)
}

// DecodeOptionFunc is the option constructor for DecodeOption.
func DecodeOptionFunc[IT any, DC instruct.DecodeContext, COPTT any](f func(o *COPTT)) *DecodeOptionImpl[IT, DC, COPTT] {
	return &DecodeOptionImpl[IT, DC, COPTT]{f}
}

// TypeDecodeOption

type TypeDecodeOptionImpl[IT any, DC instruct.DecodeContext, COPTT any] struct {
	f func(o *COPTT)
}

var _ TypeDecodeOption[any, instruct.DecodeContext, any] = (*TypeDecodeOptionImpl[any, instruct.DecodeContext, any])(nil)

func (f TypeDecodeOptionImpl[IT, DC, COPTT]) isAnyTypeOption() {}

func (f TypeDecodeOptionImpl[IT, DC, COPTT]) ApplyTypeDecodeOption(o *COPTT) {
	f.f(o)
}

// TypeDecodeOptionFunc is the option constructor for TypeDecodeOption.
func TypeDecodeOptionFunc[IT any, DC instruct.DecodeContext, COPTT any](f func(o *COPTT)) *TypeDecodeOptionImpl[IT, DC, COPTT] {
	return &TypeDecodeOptionImpl[IT, DC, COPTT]{f}
}

// DefaultOption + TypeDefaultOption

type DefaultAndTypeDefaultOptionImpl[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any] struct {
	df func(o *DOPTT)
	tf func(o *TOPTT)
}

var _ DefaultAndTypeDefaultOption[any, instruct.DecodeContext, any, any] = (*DefaultAndTypeDefaultOptionImpl[any, instruct.DecodeContext, any, any])(nil)

func (f DefaultAndTypeDefaultOptionImpl[IT, DC, DOPTT, TOPTT]) isAnyOption()     {}
func (f DefaultAndTypeDefaultOptionImpl[IT, DC, DOPTT, TOPTT]) isAnyTypeOption() {}

func (f DefaultAndTypeDefaultOptionImpl[IT, DC, DOPTT, TOPTT]) ApplyDefaultOption(o *DOPTT) {
	f.df(o)
}

func (f DefaultAndTypeDefaultOptionImpl[IT, DC, DOPTT, TOPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

// DefaultAndTypeDefaultOptionFunc is the option constructor for DefaultAndTypeDefaultOption.
func DefaultAndTypeDefaultOptionFunc[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any](df func(o *DOPTT), tf func(o *TOPTT)) *DefaultAndTypeDefaultOptionImpl[IT, DC, DOPTT, TOPTT] {
	return &DefaultAndTypeDefaultOptionImpl[IT, DC, DOPTT, TOPTT]{df, tf}
}

// DefaultOption + DecodeOption

var _ DefaultAndDecodeOption[any, instruct.DecodeContext, any, any] = (*DefaultAndDecodeOptionImpl[any, instruct.DecodeContext, any, any])(nil)

type DefaultAndDecodeOptionImpl[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any] struct {
	tf func(o *TOPTT)
	cf func(o *COPTT)
}

func (f DefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) isAnyOption() {}

func (f DefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyDefaultOption(o *TOPTT) {
	f.tf(o)
}

func (f DefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyDecodeOption(o *COPTT) {
	f.cf(o)
}

// DefaultAndDecodeOptionFunc is the option constructor for DefaultAndDecodeOption.
func DefaultAndDecodeOptionFunc[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any](tf func(o *TOPTT), cf func(o *COPTT)) *DefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT] {
	return &DefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]{tf, cf}
}

// TypeDefaultOption + TypeDecodeOption

var _ TypeDefaultAndTypeDecodeOption[any, instruct.DecodeContext, any, any] = (*TypeDefaultAndTypeDecodeOptionImpl[any, instruct.DecodeContext, any, any])(nil)

type TypeDefaultAndTypeDecodeOptionImpl[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any] struct {
	tf func(o *TOPTT)
	cf func(o *COPTT)
}

func (f TypeDefaultAndTypeDecodeOptionImpl[IT, DC, TOPTT, COPTT]) isAnyTypeOption() {}

func (f TypeDefaultAndTypeDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func (f TypeDefaultAndTypeDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyTypeDecodeOption(o *COPTT) {
	f.cf(o)
}

// TypeDefaultAndTypeDecodeOptionFunc is the option constructor for TypeDefaultAndTypeDecodeOption.
func TypeDefaultAndTypeDecodeOptionFunc[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any](tf func(o *TOPTT), cf func(o *COPTT)) *TypeDefaultAndTypeDecodeOptionImpl[IT, DC, TOPTT, COPTT] {
	return &TypeDefaultAndTypeDecodeOptionImpl[IT, DC, TOPTT, COPTT]{tf, cf}
}

// TypeDefaultOption + DecodeOption

var _ TypeDefaultAndDecodeOption[any, instruct.DecodeContext, any, any] = (*TypeDefaultAndDecodeOptionImpl[any, instruct.DecodeContext, any, any])(nil)

type TypeDefaultAndDecodeOptionImpl[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any] struct {
	tf func(o *TOPTT)
	cf func(o *COPTT)
}

func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) isAnyOption()     {}
func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) isAnyTypeOption() {}

func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyDecodeOption(o *COPTT) {
	f.cf(o)
}

// TypeDefaultAndDecodeOptionFunc is the option constructor for TypeDefaultAndDecodeOption.
func TypeDefaultAndDecodeOptionFunc[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any](tf func(o *TOPTT), cf func(o *COPTT)) *TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT] {
	return &TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]{tf, cf}
}

// DefaultOption + TypeDefaultOption + DecodeOption + TypeDecodeOption

type FullOptionImpl[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any, DCOPTT any, TCOPTT any] struct {
	df  func(o *DOPTT)
	tf  func(o *TOPTT)
	dcf func(o *DCOPTT)
	tcf func(o *TCOPTT)
}

var _ FullOption[any, instruct.DecodeContext, any, any, any, any] = (*FullOptionImpl[any, instruct.DecodeContext, any, any, any, any])(nil)

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]) isAnyOption()     {}
func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]) isAnyTypeOption() {}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]) ApplyDefaultOption(o *DOPTT) {
	f.df(o)
}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]) ApplyDecodeOption(o *DCOPTT) {
	f.dcf(o)
}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]) ApplyTypeDecodeOption(o *TCOPTT) {
	f.tcf(o)
}

// FullOptionFunc is the option constructor for FullOption.
func FullOptionFunc[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any, DCOPTT any, TCOPTT any](df func(o *DOPTT),
	tf func(o *TOPTT), dcf func(o *DCOPTT), tcf func(o *TCOPTT)) *FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT] {
	return &FullOptionImpl[IT, DC, DOPTT, TOPTT, DCOPTT, TCOPTT]{df, tf, dcf, tcf}
}
