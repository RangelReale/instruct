package options

import (
	"github.com/RangelReale/instruct"
)

// DefaultOption + TypeDefaultOption

type DefaultOptionImpl[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any] struct {
	df func(o *DOPTT)
	tf func(o *TOPTT)
}

func (f DefaultOptionImpl[IT, DC, DOPTT, TOPTT]) isOption() {}

func (f DefaultOptionImpl[IT, DC, DOPTT, TOPTT]) ApplyDefaultOption(o *DOPTT) {
	f.df(o)
}

func (f DefaultOptionImpl[IT, DC, DOPTT, TOPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func DefaultOptionFunc[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any](df func(o *DOPTT), tf func(o *TOPTT)) *DefaultOptionImpl[IT, DC, DOPTT, TOPTT] {
	return &DefaultOptionImpl[IT, DC, DOPTT, TOPTT]{df, tf}
}

// TypeDefaultOption

type TypeDefaultOptionImpl[IT any, DC instruct.DecodeContext, TOPTT any] struct {
	tf func(o *TOPTT)
}

func (f TypeDefaultOptionImpl[IT, DC, TOPTT]) isOption() {}

func (f TypeDefaultOptionImpl[IT, DC, TOPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func TypeDefaultOptionFunc[IT any, DC instruct.DecodeContext, TOPTT any](tf func(o *TOPTT)) *TypeDefaultOptionImpl[IT, DC, TOPTT] {
	return &TypeDefaultOptionImpl[IT, DC, TOPTT]{tf}
}

// DecodeOption

type DecodeOptionImpl[IT any, DC instruct.DecodeContext, COPTT any] struct {
	cf func(o *COPTT)
}

func (f DecodeOptionImpl[IT, DC, COPTT]) isOption() {}

func (f DecodeOptionImpl[IT, DC, COPTT]) ApplyDecodeOption(o *COPTT) {
	f.cf(o)
}

func DecodeOptionFunc[IT any, DC instruct.DecodeContext, COPTT any](tf func(o *COPTT)) *DecodeOptionImpl[IT, DC, COPTT] {
	return &DecodeOptionImpl[IT, DC, COPTT]{tf}
}

// TypeDefaultOption + DecodeOption

type TypeDefaultAndDecodeOptionImpl[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any] struct {
	tf func(o *TOPTT)
	cf func(o *COPTT)
}

func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) isOption() {}

func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func (f TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]) ApplyDecodeOption(o *COPTT) {
	f.cf(o)
}

func TypeDefaultAndDecodeOptionFunc[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any](tf func(o *TOPTT), cf func(o *COPTT)) *TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT] {
	return &TypeDefaultAndDecodeOptionImpl[IT, DC, TOPTT, COPTT]{tf, cf}
}

// DefaultOption + TypeDefaultOption + DecodeOption

type FullOptionImpl[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any, COPTT any] struct {
	df func(o *DOPTT)
	tf func(o *TOPTT)
	cf func(o *COPTT)
}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, COPTT]) isOption() {}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, COPTT]) ApplyDefaultOption(o *DOPTT) {
	f.df(o)
}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, COPTT]) ApplyTypeDefaultOption(o *TOPTT) {
	f.tf(o)
}

func (f FullOptionImpl[IT, DC, DOPTT, TOPTT, COPTT]) ApplyDecodeOption(o *COPTT) {
	f.cf(o)
}

func FullOptionFunc[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any, COPTT any](df func(o *DOPTT), tf func(o *TOPTT), cf func(o *COPTT)) *FullOptionImpl[IT, DC, DOPTT, TOPTT, COPTT] {
	return &FullOptionImpl[IT, DC, DOPTT, TOPTT, COPTT]{df, tf, cf}
}
