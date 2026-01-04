package tpl_test

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/pkg/tpl"
)

func TestRenderWithPanic(t *testing.T) {
	template := template.New("")
	assert.Panics(t, func() { tpl.RenderWithPanic(template, nil, "", nil) })
}
