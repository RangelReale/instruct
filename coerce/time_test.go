package coerce_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RangelReale/instruct/coerce"
	"github.com/stretchr/testify/assert"
)

// TimeTest is the struct used to build up table driven tests for time.Time.
type TimeTest struct {
	To     interface{}
	Error  error
	Expect time.Time
}

// TimeTests is a table of TimeTests.
type TimeTests map[string]*TimeTest

// Run iterates the BoolTests and runs each one.
func (tests TimeTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("time "+name, func(t *testing.T) {
			chk := assert.New(t)
			s, err := coerce.Time(test.To, time.RFC3339)
			chk.True(errors.Is(err, test.Error))
			chk.Equal(test.Expect, s)
		})
	}
}

func TestTimeUnsupported(t *testing.T) {
	tests := TimeTests{
		"unsupported": {
			To:    map[string]string{},
			Error: coerce.ErrUnsupported,
		},
	}
	tests.Run(t)
}

func TestTimeFromNil(t *testing.T) {
	chk := assert.New(t)
	var dst time.Time
	var err error
	dst, err = coerce.Time(nil, time.RFC3339)
	chk.NoError(err)
	chk.Equal(time.Time{}, dst)
}

func TestTimeFromPtr(t *testing.T) {
	n := (*time.Time)(nil)
	nn := &n
	s, _ := time.Parse(time.RFC3339, "2021-10-22T11:01:00Z")
	ps := &s
	pps := &ps
	ppps := &pps
	//
	tests := TimeTests{
		"nil": {
			To: n, Expect: time.Time{},
		},
		"*nil": {
			To: nn, Expect: time.Time{},
		},
		"ppps": {
			To: ppps, Expect: s,
		},
	}
	tests.Run(t)
}

func TestTimeFromString(t *testing.T) {
	ss := "2021-10-22T11:01:00Z"
	s, _ := time.Parse(time.RFC3339, ss)

	tests := TimeTests{
		ss: {
			To: ss, Expect: s,
		},
	}
	tests.Run(t)
}
