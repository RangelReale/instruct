package coerce_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RangelReale/instruct/coerce"
	"github.com/stretchr/testify/assert"
)

// TimeDurationTest is the struct used to build up table driven tests for time.Duration.
type TimeDurationTest struct {
	To      interface{}
	Error   error
	ErrorAs error
	Expect  time.Duration
}

// TimeDurationTests is a table of TimeDurationTests.
type TimeDurationTests map[string]*TimeDurationTest

// Run iterates the BoolTests and runs each one.
func (tests TimeDurationTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("time duration "+name, func(t *testing.T) {
			chk := assert.New(t)
			s, err := coerce.TimeDuration(test.To)
			if test.ErrorAs != nil {
				chk.ErrorAs(err, &test.ErrorAs)
			} else {
				chk.True(errors.Is(err, test.Error), "%v", err)
			}
			chk.Equal(test.Expect, s)
		})
	}
}

func TestTimeDurationUnsupported(t *testing.T) {
	tests := TimeDurationTests{
		"unsupported": {
			To:    map[string]string{},
			Error: coerce.ErrUnsupported,
		},
	}
	tests.Run(t)
}

func TestTimeDurationFromNil(t *testing.T) {
	chk := assert.New(t)
	var dst time.Duration
	var err error
	dst, err = coerce.TimeDuration(nil)
	chk.NoError(err)
	chk.Equal(time.Duration(0), dst)
}

func TestTimeDurationFromPtr(t *testing.T) {
	n := (*time.Duration)(nil)
	nn := &n
	s, _ := time.ParseDuration("5s")
	ps := &s
	pps := &ps
	ppps := &pps
	//
	tests := TimeDurationTests{
		"nil": {
			To: n, Expect: time.Duration(0),
		},
		"*nil": {
			To: nn, Expect: time.Duration(0),
		},
		"ppps": {
			To: ppps, Expect: s,
		},
	}
	tests.Run(t)
}

func TestTimeDurationFromString(t *testing.T) {
	ss := "5s"
	s, _ := time.ParseDuration(ss)

	tests := TimeDurationTests{
		ss: {
			To: ss, Expect: s,
		},
	}
	tests.Run(t)
}

func TestTimeDurationFromStringKind(t *testing.T) {
	ss := "5s"
	s, _ := time.ParseDuration(ss)

	tests := TimeDurationTests{
		"string kind empty": {
			To: S(""), Expect: time.Duration(0), ErrorAs: &time.ParseError{},
		},
		"string kind": {
			To: S(ss), Expect: s,
		},
	}
	tests.Run(t)
}
