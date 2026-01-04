package shared_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/internal/shared"
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
	assert.Equal(t, "string", ignoreError(shared.AnyToString("string")))
	assert.Equal(t, "1", ignoreError(shared.AnyToString(1)))
	assert.Equal(t, "1", ignoreError(shared.AnyToString(int64(1))))
	assert.Equal(
		t,
		"str1,str2",
		ignoreError(shared.AnyToString([]string{"str1", "str2"})),
	)
	assert.Equal(t, "1,2", ignoreError(shared.AnyToString([]int{1, 2})))
	assert.Equal(t, "1,2", ignoreError(shared.AnyToString([]int64{1, 2})))
	assert.Error(
		t,
		errors.New("undefined type"),
		ignoreValue(shared.AnyToString(Random{})),
	)
}
