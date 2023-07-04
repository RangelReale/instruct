package options

import "github.com/RangelReale/instruct"

type Option[IT any, DC instruct.DecodeContext] interface {
	isOption()
}

type DefaultOption[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any] interface {
	TypeDefaultOption[IT, DC, TOPTT]
	ApplyDefaultOption(*DOPTT)
}

type TypeDefaultOption[IT any, DC instruct.DecodeContext, OPTT any] interface {
	Option[IT, DC]
	ApplyTypeDefaultOption(*OPTT)
}

type DecodeOption[IT any, DC instruct.DecodeContext, OPTT any] interface {
	Option[IT, DC]
	ApplyDecodeOption(*OPTT)
}

type TypeDefaultAndDecodeOption[IT any, DC instruct.DecodeContext, TOPTT any, COPTT any] interface {
	TypeDefaultOption[IT, DC, TOPTT]
	DecodeOption[IT, DC, COPTT]
}

type FullOption[IT any, DC instruct.DecodeContext, DOPTT any, TOPTT any, COPTT any] interface {
	DefaultOption[IT, DC, DOPTT, TOPTT]
	DecodeOption[IT, DC, COPTT]
}
