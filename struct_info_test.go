package instruct

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStructInfo(t *testing.T) {
	type TestData1 struct {
		A string `inreq:"query"`
		B int    `inreq:"query"`
	}
	type TestData2 struct {
		X float32   `inreq:"query"`
		Y time.Time `inreq:"query"`
	}
	type TestData struct {
		First TestData1 `inreq:"recurse"`
		TestData2
	}

	type MapTestData1 struct {
		A string
		B int
	}
	type MapTestData2 struct {
		X float32
		Y time.Time
	}
	type MapTestData struct {
		First MapTestData1
		MapTestData2
	}

	tests := []struct {
		name    string
		typ     reflect.Type
		mapTags MapTags
		options DefaultOptions[*http.Request, TestDecodeContext]
		wantErr bool
	}{
		{
			name: "builds without map tags",
			typ:  reflect.TypeOf(&TestData{}),
		},
		{
			name: "builds with map tags",
			typ:  reflect.TypeOf(&MapTestData{}),
			mapTags: map[string]any{
				"First": map[string]any{
					"A": "query",
					"B": "query",
				},
				"X": "query",
				"Y": "query",
			},
		},
		{
			name: "builds with map tags not used",
			typ:  reflect.TypeOf(&MapTestData{}),
			mapTags: map[string]any{
				"First": map[string]any{
					"A": "query",
					"B": "query",
				},
				"X": "query",
				"Y": "query",
				"Z": "query",
			},
			wantErr: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.options.FieldNameMapper = DefaultFieldNameMapper
			tt.options.TagName = "inreq"

			_, err := buildStructInfo(tt.typ, tt.mapTags, tt.options)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestStructInfoModify(t *testing.T) {
	type TestData1 struct {
		A string `inreq:"query"`
		B int    `inreq:"query"`
	}
	type TestData2 struct {
		X float32   `inreq:"query"`
		Y time.Time `inreq:"query"`
	}
	type TestData struct {
		First TestData1 `inreq:"recurse"`
		TestData2
	}

	options := DefaultOptions[*http.Request, TestDecodeContext]{
		FieldNameMapper: DefaultFieldNameMapper,
		TagName:         "inreq",
	}

	typ := reflect.TypeOf(&TestData{})

	si, err := buildStructInfo(typ, nil, options)
	require.NoError(t, err)

	mapTags := map[string]any{
		"First": map[string]any{
			"A": "header",
		},
		"Y": "header",
	}

	si2, err := structInfoWithMapTags(si, mapTags, options)
	require.NoError(t, err)

	require.Equal(t, "header", si2.fieldByName("First").fieldByName("A").tag.Operation)
	require.Equal(t, "query", si2.fieldByName("First").fieldByName("B").tag.Operation)
	require.Equal(t, "query", si2.fieldByName("TestData2").fieldByName("X").tag.Operation)
	require.Equal(t, "header", si2.fieldByName("TestData2").fieldByName("Y").tag.Operation)

	// ensure the original wasn't changed
	require.Equal(t, "query", si.fieldByName("First").fieldByName("A").tag.Operation)
	require.Equal(t, "query", si.fieldByName("First").fieldByName("B").tag.Operation)
	require.Equal(t, "query", si.fieldByName("TestData2").fieldByName("X").tag.Operation)
	require.Equal(t, "query", si.fieldByName("TestData2").fieldByName("Y").tag.Operation)
}
