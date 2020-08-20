// a1 testdata.a1
package a1

import (
	"context"

	"github.com/zhaolion/gen/parser/testpkg/model"
)

// Service a1.service interface
type Service interface {
	// service interfaces

	// Foo query foo obj
	// @gen CtxFoo
	Foo(_ context.Context) (*model.Foo, error)
}

// Entry of Service
type Entry struct{}

// Foo query foo obj. just for test now
func (entry *Entry) Foo(_ context.Context) (*model.Foo, error) {
	return &model.Foo{Name: "a"}, nil
}
