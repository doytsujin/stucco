package grpc

import (
	"context"
	"fmt"

	"github.com/graphql-editor/stucco/pkg/driver"
	"github.com/graphql-editor/stucco/pkg/proto"
	protodriver "github.com/graphql-editor/stucco/pkg/proto/driver"
)

func (m *Client) UnionResolveType(input driver.UnionResolveTypeInput) (f driver.UnionResolveTypeOutput) {
	req, err := protodriver.MakeUnionResolveTypeRequest(input)
	if err == nil {
		var resp *proto.UnionResolveTypeResponse
		resp, err = m.Client.UnionResolveType(context.Background(), req)
		if err == nil {
			f = protodriver.MakeUnionResolveTypeOutput(resp)
		}
	}
	if err != nil {
		f.Error = &driver.Error{Message: err.Error()}
	}
	return
}

// UnionResolveTypeHandler union implemented by user to handle union type resolution
type UnionResolveTypeHandler interface {
	// Handle takes UnionResolveTypeInput as a type resolution input and returns
	// type name.
	Handle(driver.UnionResolveTypeInput) (string, error)
}

// UnionResolveTypeHandlerFunc is a convienience function wrapper implementing UnionResolveTypeHandler
type UnionResolveTypeHandlerFunc func(driver.UnionResolveTypeInput) (string, error)

// Handle takes UnionResolveTypeInput as a type resolution input and returns
// type name.
func (f UnionResolveTypeHandlerFunc) Handle(in driver.UnionResolveTypeInput) (string, error) {
	return f(in)
}

// UnionResolveType executes union type resolution request agains user defined function
func (m *Server) UnionResolveType(ctx context.Context, input *proto.UnionResolveTypeRequest) (f *proto.UnionResolveTypeResponse, _ error) {
	defer func() {
		if r := recover(); r != nil {
			f = &proto.UnionResolveTypeResponse{
				Error: &proto.Error{
					Msg: fmt.Sprintf("%v", r),
				},
			}
		}
	}()
	req, err := protodriver.MakeUnionResolveTypeInput(input)
	if err == nil {
		var resp string
		resp, err = m.UnionResolveTypeHandler.Handle(req)
		f = new(proto.UnionResolveTypeResponse)
		if err == nil {
			*f = protodriver.MakeUnionResolveTypeResponse(resp)
		}
	}
	if err != nil {
		f.Error = &proto.Error{Msg: err.Error()}
	}
	return
}
