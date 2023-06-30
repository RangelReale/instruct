package instruct

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"

	"golang.org/x/exp/maps"
)

func GetTestDecoderOptions() DefaultOptions[*http.Request, TestDecodeContext] {
	optns := NewDefaultOptions[*http.Request, TestDecodeContext]()
	optns.DecodeOperations[TestOperationQuery] = &TestDecodeOperationQuery{}
	optns.DecodeOperations[TestOperationHeader] = &TestDecodeOperationHeader{}
	optns.DecodeOperations[TestOperationBody] = &TestDecodeOperationBody{}
	return optns
}

func GetTestDecoderDecodeOptions(ctx TestDecodeContext) DecodeOptions[*http.Request, TestDecodeContext] {
	if ctx == nil {
		ctx = &testDecodeContext{
			sliceSplitSeparator: ",",
		}
	}

	optns := NewDecodeOptions[*http.Request, TestDecodeContext]()
	optns.Ctx = ctx
	return optns
}

type TestDecodeContext interface {
	DecodeContext
	IsBodyDecoded() bool
	DecodedBody()
	AllowReadBody() bool
	SliceSplitSeparator() string
	EnsureAllQueryUsed() bool
	EnsureAllFormUsed() bool
}

type testDecodeContext struct {
	DefaultDecodeContext
	decodedBody         bool
	allowReadBody       bool
	sliceSplitSeparator string
	ensureAllQueryUsed  bool
	ensureAllFormUsed   bool
}

func (d *testDecodeContext) IsBodyDecoded() bool {
	return d.decodedBody
}

func (d *testDecodeContext) DecodedBody() {
	d.decodedBody = true
}

func (d *testDecodeContext) AllowReadBody() bool {
	return d.allowReadBody
}

func (d *testDecodeContext) SliceSplitSeparator() string {
	return d.sliceSplitSeparator
}

func (d *testDecodeContext) EnsureAllQueryUsed() bool {
	return d.ensureAllQueryUsed
}

func (d *testDecodeContext) EnsureAllFormUsed() bool {
	return d.ensureAllFormUsed
}

const (
	TestOperationQuery  string = "query"
	TestOperationHeader        = "header"
	TestOperationBody          = "body"
)

type TestDecodeOperationQuery struct {
}

func (d *TestDecodeOperationQuery) Decode(ctx TestDecodeContext, r *http.Request, field reflect.Value,
	typ reflect.Type, tag *Tag) (bool, any, error) {
	if !r.URL.Query().Has(tag.Name) {
		return false, nil, nil
	}

	if field.Kind() == reflect.Slice {
		explode, err := tag.Options.BoolValue("explode", true)
		if err != nil {
			return false, nil, err
		}

		var value []string
		if explode {
			value = strings.Split(r.URL.Query().Get(tag.Name),
				tag.Options.Value("explodesep", ctx.SliceSplitSeparator()))
		} else {
			value = r.URL.Query()[tag.Name]
		}

		ctx.ValueUsed(TestOperationQuery, tag.Name)
		return true, value, nil
	}

	ctx.ValueUsed(TestOperationQuery, tag.Name)
	return true, r.URL.Query().Get(tag.Name), nil
}

func (d *TestDecodeOperationQuery) Validate(ctx TestDecodeContext, r *http.Request) error {
	if !ctx.EnsureAllQueryUsed() {
		return nil
	}

	queryKeys := map[string]bool{}
	for key, _ := range r.URL.Query() {
		queryKeys[key] = true
	}

	if !maps.Equal(queryKeys, ctx.GetUsedValues(TestOperationQuery)) {
		return ValuesNotUsedError{Operation: TestOperationQuery}
	}

	return nil
}

type TestDecodeOperationHeader struct {
}

func (d *TestDecodeOperationHeader) Decode(ctx TestDecodeContext, r *http.Request, field reflect.Value,
	typ reflect.Type, tag *Tag) (bool, any, error) {
	values := r.Header.Values(tag.Name)

	if len(values) == 0 {
		return false, nil, nil
	}

	if field.Kind() == reflect.Slice {
		return true, values, nil
	}
	return true, values[0], nil
}

type TestDecodeOperationBody struct {
}

func (d *TestDecodeOperationBody) Decode(ctx TestDecodeContext, r *http.Request, field reflect.Value,
	typ reflect.Type, tag *Tag) (bool, any, error) {
	if ctx.IsBodyDecoded() {
		return false, nil, fmt.Errorf("body was already decoded for type '%s'", typ.String())
	}
	fv := field
	if fv.CanAddr() {
		fv = fv.Addr()
	}
	found, err := decodeBody(ctx, r, fv.Interface(), tag)
	return found, IgnoreDecodeValue, err
}

func decodeBody(ctx TestDecodeContext, r *http.Request, data interface{}, tag *Tag) (bool, error) {
	if r.Body == nil {
		return false, nil
	}

	if !ctx.AllowReadBody() {
		return false, errors.New("body operation not allowed")
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	ctx.DecodedBody() // signal that the body was decoded

	if len(b) == 0 {
		return false, nil
	}

	var mediatype string

	if tag != nil {
		if typeStr := tag.Options.Value("type", ""); typeStr != "" {
			switch typeStr {
			case "json":
				mediatype = "application/json"
			case "xml":
				mediatype = "application/xml"
			default:
				return false, fmt.Errorf("invalid body type: '%s'", typeStr)
			}
		}
	}

	if mediatype == "" {
		mediatype, _, err = mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			return false, fmt.Errorf("error detecting body content type: %w", err)
		}
	}

	switch mediatype {
	case "application/json":
		err := json.Unmarshal(b, &data)
		if err != nil {
			return true, fmt.Errorf("error parsing JSON body: %w", err)
		}
		return true, nil
	case "text/xml", "application/xml":
		err := xml.Unmarshal(b, &data)
		if err != nil {
			return true, fmt.Errorf("error parsing XML body: %w", err)
		}
		return true, nil
	}

	return false, nil
}
