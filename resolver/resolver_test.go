package resolver

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_resolvesValues(t *testing.T) {
	resolver := NewResolver()

	tests := []struct {
		name    string
		input   interface{}
		value   []string
		want    interface{}
		wantErr bool
	}{
		{name: "resolve []string", input: []string{}, value: []string{"test"}, want: []string{"test"}, wantErr: false},
		{name: "failed unsupported type", input: []struct{}{}, value: []string{"trick"}, want: []struct{}(nil), wantErr: true},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			f := reflect.New(reflect.TypeOf(tt.input)).Elem()
			err := resolver.Resolve(f, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, f.Interface())
		})
	}
}

func Test_resolveValue(t *testing.T) {
	var ptrInput *bool
	b := true
	var structInput *struct{}
	resolver := NewResolver()
	tests := []struct {
		name    string
		input   interface{}
		value   string
		want    interface{}
		wantErr bool
	}{
		{name: "resolve string", input: string(""), value: "test", want: "test", wantErr: false},
		{name: "resolve pointer", input: ptrInput, value: "true", want: &b, wantErr: false},
		{name: "failed unsupported type", input: struct{}{}, value: "trick", want: struct{}{}, wantErr: true},
		{name: "failed unsupported pointertype", input: structInput, value: "trick", want: structInput, wantErr: true},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			f := reflect.New(reflect.TypeOf(tt.input)).Elem()
			err := resolver.Resolve(f, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, f.Interface())
		})
	}
}

func Test_resolve(t *testing.T) {
	resolver := &DefaultValueResolver{
		CustomTypes: []TypeValueResolver{
			NewValueResolverTime(time.RFC3339),
			&ValueResolverTimeDuration{},
		},
	}

	t1, _ := time.Parse(time.RFC3339, "2021-10-22T11:01:00Z")
	tests := []struct {
		name    string
		input   interface{}
		value   string
		want    interface{}
		wantErr bool
	}{
		{name: "resolve string", input: string(""), value: "test", want: "test", wantErr: false},
		{name: "resolve bool", input: bool(false), value: "true", want: true, wantErr: false},
		{name: "resolve failed bool", input: bool(false), value: "trick", want: bool(false), wantErr: true},
		{name: "resolve time", input: time.Time{}, value: "2021-10-22T11:01:00Z", want: t1, wantErr: false},
		{name: "resolve failed time", input: time.Time{}, value: "trick", want: time.Time{}, wantErr: true},
		{name: "resolve duration", input: time.Duration(0), value: "5s", want: 5 * time.Second, wantErr: false},
		{name: "resolve failed duration", input: time.Duration(0), value: "trick", want: time.Duration(0), wantErr: true},
		{name: "resolve int", input: int(0), value: "5", want: int(5), wantErr: false},
		{name: "resolve failed int", input: int(0), value: "trick", want: int(0), wantErr: true},
		{name: "resolve int64", input: int64(0), value: "5", want: int64(5), wantErr: false},
		{name: "resolve failed int64", input: int64(0), value: "trick", want: int64(0), wantErr: true},
		{name: "resolve int32", input: int32(0), value: "5", want: int32(5), wantErr: false},
		{name: "resolve failed int32", input: int32(0), value: "trick", want: int32(0), wantErr: true},
		{name: "resolve int16", input: int16(0), value: "5", want: int16(5), wantErr: false},
		{name: "resolve failed int16", input: int16(0), value: "trick", want: int16(0), wantErr: true},
		{name: "resolve int8", input: int8(0), value: "5", want: int8(5), wantErr: false},
		{name: "resolve failed int8", input: int8(0), value: "trick", want: int8(0), wantErr: true},
		{name: "resolve float64", input: float64(0), value: "5.5", want: float64(5.5), wantErr: false},
		{name: "resolve failed float64", input: float64(0), value: "trick", want: float64(0), wantErr: true},
		{name: "resolve float32", input: float32(0), value: "5.5", want: float32(5.5), wantErr: false},
		{name: "resolve failed float32", input: float32(0), value: "trick", want: float32(0), wantErr: true},
		{name: "resolve uint", input: uint(0), value: "5", want: uint(5), wantErr: false},
		{name: "resolve failed uint", input: uint(0), value: "trick", want: uint(0), wantErr: true},
		{name: "resolve uint64", input: uint64(0), value: "5", want: uint64(5), wantErr: false},
		{name: "resolve failed uint64", input: uint64(0), value: "trick", want: uint64(0), wantErr: true},
		{name: "resolve uint32", input: uint32(0), value: "5", want: uint32(5), wantErr: false},
		{name: "resolve failed uint32", input: uint32(0), value: "trick", want: uint32(0), wantErr: true},
		{name: "resolve uint16", input: uint16(0), value: "5", want: uint16(5), wantErr: false},
		{name: "resolve failed uint16", input: uint16(0), value: "trick", want: uint16(0), wantErr: true},
		{name: "resolve uint8", input: uint8(0), value: "5", want: uint8(5), wantErr: false},
		{name: "resolve failed uint8", input: uint8(0), value: "trick", want: uint8(0), wantErr: true},
		{name: "failed unsupported type", input: []struct{}{}, value: "trick", want: nil, wantErr: true},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			target := reflect.Indirect(reflect.New(reflect.TypeOf(tt.input)))
			target.Set(reflect.ValueOf(tt.input))
			err := resolver.ResolveValue(target, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				require.Equal(t, tt.want, target.Interface())
			}
		})
	}
}

func Test_resolve_textUnmarshaller(t *testing.T) {
	resolver := &DefaultValueResolver{
		CustomTypesReflect: []TypeValueResolverReflect{
			&ValueResolverReflectTextUnmarshaler{},
		},
	}

	t1, _ := time.Parse(time.RFC3339, "2021-10-22T11:01:00Z")

	target := reflect.Indirect(reflect.New(reflect.TypeOf(time.Time{})))
	err := resolver.ResolveValue(target, "2021-10-22T11:01:00Z")
	require.NoError(t, err)
	require.Equal(t, t1, target.Interface())
}
