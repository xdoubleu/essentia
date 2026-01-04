//nolint:exhaustruct //on purpose
package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/validate"
)

type TestStruct struct {
	StrVal    string  `json:"strVal"`
	IntVal    int     `json:"intVal"`
	Int64Val  int64   `json:"int64Val"`
	TzVal     string  `json:"tzVal"`
	OptStrVal *string `json:"optStrVal"`
}

func (ts *TestStruct) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "strVal", ts.StrVal, validate.IsNotEmpty)

	validate.Check(v, "intVal", ts.IntVal, validate.IsGreaterThan(-1))
	validate.Check(v, "intVal", ts.IntVal, validate.IsGreaterThanOrEqual(1))
	validate.Check(v, "intVal", ts.IntVal, validate.IsLesserThan(4))
	validate.Check(v, "intVal", ts.IntVal, validate.IsLesserThanOrEqual(2))

	validate.Check(v, "int64Val", ts.Int64Val, validate.IsGreaterThan(int64(-1)))
	validate.Check(
		v,
		"int64Val",
		ts.Int64Val,
		validate.IsGreaterThanOrEqual(int64(1)),
	)
	validate.Check(v, "int64Val", ts.Int64Val, validate.IsLesserThan(int64(4)))
	validate.Check(
		v,
		"int64Val",
		ts.Int64Val,
		validate.IsLesserThanOrEqual(int64(2)),
	)

	validate.Check(v, "tzVal", ts.TzVal, validate.IsValidTimeZone)

	validate.CheckOptional(
		v,
		"optStrVal",
		ts.OptStrVal,
		validate.IsInSlice([]string{"allowed"}),
	)

	return v.Valid(), v.Errors()
}

func TestAllOk(t *testing.T) {
	val := "allowed"

	ts := TestStruct{
		StrVal:    "hello",
		IntVal:    1,
		Int64Val:  1,
		TzVal:     "Europe/Brussels",
		OptStrVal: &val,
	}

	valid, errors := ts.Validate()
	assert.True(t, valid)
	assert.Equal(t, 0, len(errors))
}

func TestIsEmpty(t *testing.T) {
	ts := TestStruct{
		StrVal:   "",
		IntVal:   1,
		Int64Val: 1,
		TzVal:    "Europe/Brussels",
	}

	expectedErrors := map[string]string{
		"strVal": "must be provided",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}

func TestIsNotGT(t *testing.T) {
	ts := TestStruct{
		StrVal:   "hello",
		IntVal:   -1,
		Int64Val: -1,
		TzVal:    "Europe/Brussels",
	}

	expectedErrors := map[string]string{
		"intVal":   "must be greater than -1",
		"int64Val": "must be greater than -1",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}

func TestIsNotGTE(t *testing.T) {
	ts := TestStruct{
		StrVal:   "hello",
		IntVal:   0,
		Int64Val: 0,
		TzVal:    "Europe/Brussels",
	}

	expectedErrors := map[string]string{
		"intVal":   "must be greater than or equal to 1",
		"int64Val": "must be greater than or equal to 1",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}

func TestIsNotLT(t *testing.T) {
	ts := TestStruct{
		StrVal:   "hello",
		IntVal:   4,
		Int64Val: 4,
		TzVal:    "Europe/Brussels",
	}

	expectedErrors := map[string]string{
		"intVal":   "must be lesser than 4",
		"int64Val": "must be lesser than 4",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}

func TestIsNotLTE(t *testing.T) {
	ts := TestStruct{
		StrVal:   "hello",
		IntVal:   3,
		Int64Val: 3,
		TzVal:    "Europe/Brussels",
	}

	expectedErrors := map[string]string{
		"intVal":   "must be lesser than or equal to 2",
		"int64Val": "must be lesser than or equal to 2",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}

func TestIsNotValidTz(t *testing.T) {
	ts := TestStruct{
		StrVal:   "hello",
		IntVal:   1,
		Int64Val: 1,
		TzVal:    "AAAAAAAAAAAAAAAAAAAAAA",
	}

	expectedErrors := map[string]string{
		"tzVal": "must be a valid IANA value",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}

func TestIsNotInSlice(t *testing.T) {
	val := "notallowed"
	ts := TestStruct{
		StrVal:    "hello",
		IntVal:    1,
		Int64Val:  1,
		TzVal:     "Europe/Brussels",
		OptStrVal: &val,
	}

	expectedErrors := map[string]string{
		"optStrVal": "must be a valid value",
	}

	valid, errors := ts.Validate()
	assert.False(t, valid)
	assert.Equal(t, expectedErrors, errors)
}
