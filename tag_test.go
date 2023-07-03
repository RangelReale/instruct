package instruct

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTag(t *testing.T) {
	defOpt := GetTestDecoderOptions()

	tests := []struct {
		name              string
		fieldName         string
		tagValue          string
		expectedError     bool
		expectedName      string
		expectedOperation string
		expectedRequired  bool
		expectedOptions   map[string]string
		expectedSOWhen    string
		expectedSORecurse bool
	}{
		{
			name:              "with operation",
			fieldName:         "Val",
			tagValue:          "header",
			expectedName:      "val",
			expectedOperation: "header",
			expectedRequired:  true,
			expectedOptions:   map[string]string{},
		},
		{
			name:              "with options",
			fieldName:         "Val",
			tagValue:          "header,a=1,b=2",
			expectedName:      "val",
			expectedOperation: "header",
			expectedRequired:  true,
			expectedOptions: map[string]string{
				"a": "1",
				"b": "2",
			},
		},
		{
			name:              "not required",
			fieldName:         "Val",
			tagValue:          "header,required=false",
			expectedName:      "val",
			expectedOperation: "header",
			expectedRequired:  false,
			expectedOptions:   map[string]string{},
		},
		{
			name:          "required invalid",
			fieldName:     "Val",
			tagValue:      "header,required=invalid_value",
			expectedError: true,
		},
		{
			name:          "empty operation error",
			fieldName:     "Val",
			tagValue:      ",a=1,b=2",
			expectedError: true,
		},
		{
			name:              "empty name uses field name",
			fieldName:         "Val",
			tagValue:          "header,name=,a=1,b=2",
			expectedName:      "val",
			expectedOperation: "header",
			expectedRequired:  true,
			expectedOptions: map[string]string{
				"a": "1",
				"b": "2",
			},
		},
		{
			name:          "invalid options",
			fieldName:     "Val",
			tagValue:      "header,a,b",
			expectedError: true,
		},
		{
			name:              "ignore empty options",
			fieldName:         "Val",
			tagValue:          "header,,",
			expectedName:      "val",
			expectedOperation: "header",
			expectedRequired:  true,
			expectedOptions:   map[string]string{},
		},
		{
			name:              "with SO options",
			fieldName:         "Val",
			tagValue:          "header,so_when=before,so_recurse=true",
			expectedName:      "val",
			expectedOperation: "header",
			expectedRequired:  true,
			expectedOptions:   map[string]string{},
			expectedSOWhen:    "before",
			expectedSORecurse: true,
		},
		{
			name:          "with SO option error",
			fieldName:     "Val",
			tagValue:      "header,so_when=invalid_value",
			expectedError: true,
		},
		{
			name:          "with SO invalid option",
			fieldName:     "Val",
			tagValue:      "header,so_invalid=invalid_value",
			expectedError: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			tag, err := parseTags[*http.Request, TestDecodeContext](tt.fieldName, tt.tagValue, &defOpt)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedName, tag.Name)
				require.Equal(t, tt.expectedOperation, tag.Operation)
				require.Equal(t, tt.expectedRequired, tag.Required)
				require.Equal(t, tt.expectedOptions, tag.Options.options)
				require.Equal(t, tt.expectedSOWhen, tag.SOWhen)
				require.Equal(t, tt.expectedSORecurse, tag.SORecurse)
			}
		})
	}

}
