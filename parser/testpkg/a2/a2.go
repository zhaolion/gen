package a2

import (
	"context"

	"github.com/zhaolion/gen/parser/testpkg/a1"
	"github.com/zhaolion/gen/parser/testpkg/model"
)

// Entry of a2.Service. wrapper of a1.Entry
type Entry struct {
	A1 a1.Entry
}

// Foo query foo obj. just for test now
func (entry *Entry) Foo(ctx context.Context) (*model.Foo, error) {
	return entry.A1.Foo(ctx)
}

// Bar query foo obj. just for test a2 now
func (entry *Entry) Bar(ctx context.Context) (*model.Foo, error) {
	return entry.A1.Foo(ctx)
}
