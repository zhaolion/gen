package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	builder := New()
	builder.SetDebugLevel()
	builder.AddDirRecursive("./testpkg")

	universe, err := builder.FindTypes()
	if err != nil {
		t.Fatalf("invalid FindTypes err: %+v", err)
	}
	assert.NotEmpty(t, universe)
}
