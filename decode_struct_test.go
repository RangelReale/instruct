package instruct

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeStructInvalidType(t *testing.T) {
	type DataType struct {
		Val string `instruct:"header"`
	}
	type DataType2 struct {
		Val string `instruct:"query"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())

	si, err := dec.structInfoFromType(reflect.TypeOf(DataType2{}))
	require.NoError(t, err)

	err = dec.decodeInputFromStructInfo(r, si, &data, GetTestDecoderDecodeOptions(nil))
	require.Error(t, err)
}
