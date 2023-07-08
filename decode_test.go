package instruct

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/RangelReale/instruct/types"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	type DTestEmbed struct {
		H string `instruct:"header"`
		Q string `instruct:"query"`
	}

	type DTest1 struct {
		Q int `instruct:"query,name=Q1"`
	}

	type DTestBody struct {
		F1 string
		F2 int
	}

	type DTest struct {
		DTestEmbed
		T1 DTest1    `instruct:"recurse"`
		TB DTestBody `instruct:"body"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"F1":"ValueF1","F2":99}`))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("H", "ValueH")
	q := r.URL.Query()
	q.Add("q", "ValueQ")
	q.Add("Q1", "66")
	r.URL.RawQuery = q.Encode()

	data := &DTest{}
	want := &DTest{
		DTestEmbed: DTestEmbed{
			H: "ValueH",
			Q: "ValueQ",
		},
		T1: DTest1{
			Q: 66,
		},
		TB: DTestBody{
			F1: "ValueF1",
			F2: 99,
		},
	}

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, data, GetTestDecoderDecodeOptions(&testDecodeContext{
		allowReadBody: true,
	}))
	require.NoError(t, err)
	require.Equal(t, want, data)
}

func TestDecodeEmbed(t *testing.T) {
	type EmbedTestInner struct {
		Val string `instruct:"header"`
	}

	type EmbedTest struct {
		EmbedTestInner
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	data := &EmbedTest{}
	want := &EmbedTest{
		EmbedTestInner{Val: "x1"},
	}

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, want, data)
}

func TestDecodeUsingCache(t *testing.T) {
	type DataType struct {
		Val string `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/?val=x2", nil)
	r.Header.Set("val", "x1")

	var data DataType

	defOpt := GetTestDecoderOptions()
	defOpt.StructInfoCache(true)
	dec := NewDecoder[*http.Request, TestDecodeContext](defOpt)
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)

	// this resets the cache
	dec.options.DefaultMapTagsSet(reflect.TypeOf(DataType{}), MapTags{
		"Val": "query",
	})
	err = dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "x2", data.Val)
}

func TestDecodeNonPointer(t *testing.T) {
	type DataType struct {
		Val string `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, data, GetTestDecoderDecodeOptions(nil))
	var target *types.InvalidDecodeError
	require.ErrorAs(t, err, &target)
}

func TestDecodeNoContext(t *testing.T) {
	type DataType struct {
		Val string `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	decOpt := GetTestDecoderDecodeOptions(nil)
	decOpt.Ctx = nil
	err := dec.Decode(r, &data, decOpt)
	require.Error(t, err)
}

func TestDecodePointerField(t *testing.T) {
	type SData struct {
		Val    string `instruct:"header"`
		IntVal *int32 `instruct:"header"`
	}

	type DataType struct {
		Val    *string  `instruct:"header"`
		S      *SData   `instruct:"recurse"`
		S2     ***SData `instruct:"recurse"`
		IntVal *int32   `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("intval", "92")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.NotNil(t, data.Val)
	require.NotNil(t, data.IntVal)
	require.NotNil(t, data.S)
	require.NotNil(t, data.S.IntVal)
	require.Equal(t, "x1", *data.Val)
	require.Equal(t, int32(92), *data.IntVal)
	require.Equal(t, "x1", data.S.Val)
	require.Equal(t, int32(92), *data.S.IntVal)
}

func TestDecodePointerPointerField(t *testing.T) {
	type DataType struct {
		Val **string `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.NotNil(t, data.Val)
	require.NotNil(t, *data.Val)
	require.Equal(t, "x1", **data.Val)
}

func TestDecodeSliceField(t *testing.T) {
	type DataType struct {
		Val []int32 `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Add("val", "12")
	r.Header.Add("val", "13")
	r.Header.Add("val", "15")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, []int32{12, 13, 15}, data.Val)
}

func TestDecodeSlicePointerField(t *testing.T) {
	type DataType struct {
		Val []*int `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Add("val", "12")
	r.Header.Add("val", "13")
	r.Header.Add("val", "15")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)

	v1 := 12
	v2 := 13
	v3 := 15

	require.Equal(t, []*int{&v1, &v2, &v3}, data.Val)
}

func TestDecodeRequiredError(t *testing.T) {
	type DataType struct {
		Val string `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("other-val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	var reqErr types.RequiredError
	require.ErrorAs(t, err, &reqErr)
	require.Equal(t, TestOperationHeader, reqErr.Operation)
	require.Equal(t, "Val", reqErr.FieldName)
	require.Equal(t, "val", reqErr.TagName)
}

func TestDecodeStructOption(t *testing.T) {
	type DataType struct {
		_   StructOption `instruct:"body,type=json"`
		Val string
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"Val": "14"}`))

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(&testDecodeContext{
		allowReadBody: true,
	}))
	require.NoError(t, err)
	require.Equal(t, "14", data.Val)
}

func TestDecodeMapStructOption(t *testing.T) {
	type DataType struct {
		Val string
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"Val": "14"}`))

	var data DataType

	defOpt := GetTestDecoderOptions()
	defOpt.DefaultMapTagsSet(reflect.TypeOf(DataType{}), MapTags{
		StructOptionMapTag: "body,type=json",
	})
	dec := NewDecoder[*http.Request, TestDecodeContext](defOpt)
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(&testDecodeContext{
		allowReadBody: true,
	}))
	require.NoError(t, err)
	require.Equal(t, "14", data.Val)
}

func TestDecodeStructOptionRequiredError(t *testing.T) {
	type DataType struct {
		_ StructOption `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("other-val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	var reqErr types.RequiredError
	require.ErrorAs(t, err, &reqErr)
	require.Equal(t, TestOperationHeader, reqErr.Operation)
	require.Equal(t, "instruct.DataType", reqErr.FieldName)
	require.Equal(t, "_", reqErr.TagName)
}

func TestDecodeStructOptionInvalidType(t *testing.T) {
	type DataType struct {
		_ StructOption `instruct:"header,name=val"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	var opErr types.OperationNotSupportedError
	require.ErrorAs(t, err, &opErr)
	require.Equal(t, TestOperationHeader, opErr.Operation)
	require.Equal(t, "instruct.DataType", opErr.FieldName)
}

func TestDecodeStructOptionPriority(t *testing.T) {
	type Inner struct {
		_       StructOption `instruct:"body,type=xml"`
		XMLName xml.Name     `instruct:"-" xml:"Inner"`
		Val     string
	}

	// struct tag have priority over StructOption
	type DataType struct {
		I Inner `instruct:"body,type=json"`
	}

	type DataType2 struct {
		I Inner
	}

	// map tag have priority over all others
	type DataType3 struct {
		I Inner `instruct:"body,type=json"`
	}

	defOpt := GetTestDecoderOptions()
	defOpt.DefaultMapTagsSet(reflect.TypeOf(DataType3{}), MapTags{
		"I": MapTags{
			"Val": "header",
		},
	})
	dec := NewDecoder[*http.Request, TestDecodeContext](defOpt)

	// use struct tag
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"Val": "14"}`))

	var data DataType

	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(&testDecodeContext{
		allowReadBody: true,
	}))
	require.NoError(t, err)
	require.Equal(t, "14", data.I.Val)

	// use struct option
	r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`<Inner><Val>15</Val></Inner>`))

	var data2 DataType2

	err = dec.Decode(r, &data2, GetTestDecoderDecodeOptions(&testDecodeContext{
		allowReadBody: true,
	}))
	require.NoError(t, err)
	require.Equal(t, "15", data2.I.Val)

	// use map tag
	r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`<Inner><Val>15</Val></Inner>`))
	r.Header.Set("val", "90")

	var data3 DataType3

	err = dec.Decode(r, &data3, GetTestDecoderDecodeOptions(&testDecodeContext{
		allowReadBody: true,
	}))
	require.NoError(t, err)
	require.Equal(t, "90", data3.I.Val)
}

func TestDecodeMapTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	var data DataType

	defOpt := GetTestDecoderOptions()
	defOpt.DefaultMapTagsSet(reflect.TypeOf(DataType{}), map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	})

	dec := NewDecoder[*http.Request, TestDecodeContext](defOpt)
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
}

func TestDecodeMapTagsOverrideStructTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/?val=x1", nil)
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string `instruct:"query"`
		X   struct {
			X1 string `instruct:"query"`
		}
	}

	var data DataType

	defOpt := GetTestDecoderOptions()
	defOpt.DefaultMapTagsSet(reflect.TypeOf(DataType{}), map[string]any{
		"X": map[string]any{
			"X1": "header",
		},
	})

	dec := NewDecoder[*http.Request, TestDecodeContext](defOpt)
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
	require.Equal(t, "x2", data.X.X1)
}

func TestDecodeMapTagsDecodeOverride(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/?val=x1", nil)
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string `instruct:"query"`
		X   struct {
			X1 string `instruct:"query"`
		} `instruct:"recurse"`
	}

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	decOpt := GetTestDecoderDecodeOptions(nil)
	decOpt.MapTags = map[string]any{
		"X": map[string]any{
			"X1": "header",
		},
	}

	err := dec.Decode(r, &data, decOpt)
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
	require.Equal(t, "x2", data.X.X1)
}

func TestDecodeMapTagsDecodeAsDefault(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	decOpt := GetTestDecoderDecodeOptions(nil)
	decOpt.UseDecodeMapTagsAsDefault = true
	decOpt.MapTags = map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	}
	err := dec.Decode(r, &data, decOpt)
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
}

func TestDecodeMapTagsNoDecodeAsDefault(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	decOpt := GetTestDecoderDecodeOptions(nil)
	decOpt.UseDecodeMapTagsAsDefault = false
	decOpt.MapTags = map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	}
	err := dec.Decode(r, &data, decOpt)
	require.Error(t, err)
}

func TestDecodePointerValue(t *testing.T) {
	type DataType struct {
		Val string `instruct:"manual"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	var data DataType

	mval := 45

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptionsWithManual(map[string]any{
		"val": &mval,
	}))
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
}

func TestDecodeSliceToNonSliceError(t *testing.T) {
	type DataType struct {
		Val string `instruct:"manual"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptionsWithManual(map[string]any{
		"val": []string{"A", "B"},
	}))
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.ErrorIs(t, err, types.ErrCoerceUnsupported)
}

func TestDecodeMapFieldNotSupported(t *testing.T) {
	type DataType struct {
		Val map[string]bool `instruct:"manual"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptionsWithManual(map[string]any{
		"val": "x",
	}))
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.ErrorIs(t, err, types.ErrCoerceUnknown)
}

func TestDecodeEnsureUsed(t *testing.T) {
	type DataType struct {
		Val string `instruct:"query"`
	}

	r := httptest.NewRequest(http.MethodPost, "/?val=x1&y=2", nil)

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(&testDecodeContext{
		ensureAllQueryUsed: true,
	}))
	var uerr types.ValuesNotUsedError
	require.ErrorAs(t, err, &uerr)
	require.Equal(t, "query", uerr.Operation)
}

func TestDecodeUnknownOperation(t *testing.T) {
	type DataType struct {
		Val string `instruct:"invalid_operation"`
	}

	r := httptest.NewRequest(http.MethodPost, "/?val=x1&y=2", nil)

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, &data, GetTestDecoderDecodeOptions(nil))
	require.Error(t, err)
}
