package instruct

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	type DTestEmbed struct {
		H string `instruct:"header"`
		Q string `instruct:"query"`
	}

	type DTest1 struct {
		Q string `instruct:"query,name=Q1"`
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
	q.Add("Q1", "ValueQ1")
	r.URL.RawQuery = q.Encode()

	data := &DTest{}
	want := &DTest{
		DTestEmbed: DTestEmbed{
			H: "ValueH",
			Q: "ValueQ",
		},
		T1: DTest1{
			Q: "ValueQ1",
		},
		TB: DTestBody{
			F1: "ValueF1",
			F2: 99,
		},
	}

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, data, GetTestDecoderDecodeOptions(&testDecodeContext{
		sliceSplitSeparator: ",",
		allowReadBody:       true,
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

func TestDecodeNonPointer(t *testing.T) {
	type DataType struct {
		Val string `instruct:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	err := dec.Decode(r, data, GetTestDecoderDecodeOptions(nil))
	var target *InvalidDecodeError
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
	err := dec.Decode(r, data, decOpt)
	require.Error(t, err)
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

	dec := NewDecoder[*http.Request, TestDecodeContext](GetTestDecoderOptions())
	decOpt := GetTestDecoderDecodeOptions(nil)
	decOpt.UseDecodeMapTagsAsDefault = true
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
