package instruct

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeType(t *testing.T) {
	type DataType struct {
		Val string `instruct:"query"`
	}

	d := NewTypeDecoder[*http.Request, TestDecodeContext, DataType](GetTestTypeDecoderOptions())

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	q := r.URL.Query()
	q.Add("val", "v1")
	r.URL.RawQuery = q.Encode()

	v, err := d.Decode(r, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "v1", v.Val)
}

func TestDecodeTypeMapTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	defOpt := GetTestTypeDecoderOptions()
	defOpt.DefaultOptions.defaultMapTags.Set(reflect.TypeOf(DataType{}), map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	})

	d := NewTypeDecoder[*http.Request, TestDecodeContext, DataType](defOpt)
	data, err := d.Decode(r, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
}

func TestDecodeTypeMapTagsOverrideStructTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/?val=x1", nil)
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string `instruct:"query"`
		X   struct {
			X1 string `instruct:"query"`
		}
	}

	defOpt := GetTestTypeDecoderOptions()
	defOpt.defaultMapTags.Set(reflect.TypeOf(DataType{}), map[string]any{
		"X": map[string]any{
			"X1": "header",
		},
	})

	d := NewTypeDecoder[*http.Request, TestDecodeContext, DataType](defOpt)
	data, err := d.Decode(r, GetTestDecoderDecodeOptions(nil))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
	require.Equal(t, "x2", data.X.X1)
}
