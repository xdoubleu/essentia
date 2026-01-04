package helpers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/internal/helpers"
)

type Random struct {
}

func ignoreError(value string, _ error) string {
	return value
}

func ignoreValue(_ string, err error) error {
	return err
}

func TestAnyToString(t *testing.T) {
	assert.Equal(t, "string", ignoreError(helpers.AnyToString("string")))
	assert.Equal(t, "1", ignoreError(helpers.AnyToString(1)))
	assert.Equal(t, "1", ignoreError(helpers.AnyToString(int64(1))))
	assert.Equal(
		t,
		"str1,str2",
		ignoreError(helpers.AnyToString([]string{"str1", "str2"})),
	)
	assert.Equal(t, "1,2", ignoreError(helpers.AnyToString([]int{1, 2})))
	assert.Equal(t, "1,2", ignoreError(helpers.AnyToString([]int64{1, 2})))
	assert.Error(
		t,
		errors.New("undefined type"),
		ignoreValue(helpers.AnyToString(Random{})),
	)
}
