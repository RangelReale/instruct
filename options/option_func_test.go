package options

import (
	"testing"

	"github.com/RangelReale/instruct"
	"github.com/stretchr/testify/require"
)

type testOptions struct {
	value1 int
	value2 int
	value3 int
	value4 int
}

func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		option any
		error  bool
		values [4]int
	}{
		{
			name: "default option",
			option: func() any {
				return DefaultOptionFunc[any, instruct.DecodeContext, testOptions](func(o *testOptions) {
					o.value1 = 10
				})
			}(),
			values: [4]int{10, 0, 0, 0},
		},
		{
			name: "type default option",
			option: func() any {
				return TypeDefaultOptionFunc[any, instruct.DecodeContext, testOptions](func(o *testOptions) {
					o.value1 = 15
				})
			}(),
			values: [4]int{15, 0, 0, 0},
		},
		{
			name: "decode option",
			option: func() any {
				return DecodeOptionFunc[any, instruct.DecodeContext, testOptions](func(o *testOptions) {
					o.value1 = 13
				})
			}(),
			values: [4]int{13, 0, 0, 0},
		},
		{
			name: "type decode option",
			option: func() any {
				return TypeDecodeOptionFunc[any, instruct.DecodeContext, testOptions](func(o *testOptions) {
					o.value1 = 44
				})
			}(),
			values: [4]int{44, 0, 0, 0},
		},
		{
			name: "default and type default option",
			option: func() any {
				return DefaultAndTypeDefaultOptionFunc[any, instruct.DecodeContext, testOptions, testOptions](func(o *testOptions) {
					o.value1 = 66
				}, func(o *testOptions) {
					o.value2 = 68
				})
			}(),
			values: [4]int{66, 68, 0, 0},
		},
		{
			name: "default and decode option",
			option: func() any {
				return DefaultAndDecodeOptionFunc[any, instruct.DecodeContext, testOptions, testOptions](func(o *testOptions) {
					o.value1 = 89
				}, func(o *testOptions) {
					o.value2 = 31
				})
			}(),
			values: [4]int{89, 31, 0, 0},
		},
		{
			name: "type default and type decode option",
			option: func() any {
				return TypeDefaultAndTypeDecodeOptionFunc[any, instruct.DecodeContext, testOptions, testOptions](func(o *testOptions) {
					o.value1 = 71
				}, func(o *testOptions) {
					o.value2 = 29
				})
			}(),
			values: [4]int{71, 29, 0, 0},
		},
		{
			name: "type default and decode option",
			option: func() any {
				return TypeDefaultAndDecodeOptionFunc[any, instruct.DecodeContext, testOptions, testOptions](func(o *testOptions) {
					o.value1 = 19
				}, func(o *testOptions) {
					o.value2 = 11
				})
			}(),
			values: [4]int{19, 11, 0, 0},
		},
		{
			name: "full options",
			option: func() any {
				return FullOptionFunc[any, instruct.DecodeContext, testOptions, testOptions](func(o *testOptions) {
					o.value1 = 78
				}, func(o *testOptions) {
					o.value2 = 79
				}, func(o *testOptions) {
					o.value3 = 80
				}, func(o *testOptions) {
					o.value4 = 81
				})
			}(),
			values: [4]int{78, 79, 80, 81},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			var options testOptions
			if opt, ok := tt.option.(DefaultOption[any, instruct.DecodeContext, testOptions]); ok {
				opt.ApplyDefaultOption(&options)
			}
			if opt, ok := tt.option.(TypeDefaultOption[any, instruct.DecodeContext, testOptions]); ok {
				opt.ApplyTypeDefaultOption(&options)
			}
			if opt, ok := tt.option.(DecodeOption[any, instruct.DecodeContext, testOptions]); ok {
				opt.ApplyDecodeOption(&options)
			}
			if opt, ok := tt.option.(TypeDecodeOption[any, instruct.DecodeContext, testOptions]); ok {
				opt.ApplyTypeDecodeOption(&options)
			}

			require.Equal(t, tt.values[0], options.value1)
			require.Equal(t, tt.values[1], options.value2)
			require.Equal(t, tt.values[2], options.value3)
			require.Equal(t, tt.values[3], options.value4)
		})
	}
}
