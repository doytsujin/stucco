package protohttp

import (
	"io/ioutil"
	"net/http"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/graphql-editor/stucco/pkg/driver"
	"github.com/graphql-editor/stucco/pkg/proto"
	protodriver "github.com/graphql-editor/stucco/pkg/proto/driver"
)

// InterfaceResolveType over http
func (c *Client) InterfaceResolveType(input driver.InterfaceResolveTypeInput) driver.InterfaceResolveTypeOutput {
	var out driver.InterfaceResolveTypeOutput
	req, err := protodriver.MakeInterfaceResolveTypeRequest(input)
	if err == nil {
		resp := new(proto.InterfaceResolveTypeResponse)
		if err = c.do(message{
			contentType: interfaceResolveTypeRequestMessage,
			proto:       req,
		}, message{
			contentType: interfaceResolveTypeResponseMessage,
			proto:       resp,
		}); err == nil {
			out = protodriver.MakeInterfaceResolveTypeOutput(resp)
		}
	}
	if err != nil {
		out.Error = &driver.Error{
			Message: err.Error(),
		}
	}
	return out
}

func (h *Handler) interfaceResolveType(req *http.Request) *proto.InterfaceResolveTypeResponse {
	resp := new(proto.InterfaceResolveTypeResponse)
	protoReq := new(proto.InterfaceResolveTypeRequest)
	var err error
	var b []byte
	if b, err = ioutil.ReadAll(req.Body); err == nil {
		defer req.Body.Close()
		if err = protobuf.Unmarshal(b, protoReq); err == nil {
			var in driver.InterfaceResolveTypeInput
			in, err = protodriver.MakeInterfaceResolveTypeInput(protoReq)
			if err == nil {
				var interfaceType string
				interfaceType, err = h.InterfaceResolveType(in)
				if err == nil {
					*resp = protodriver.MakeInterfaceResolveTypeResponse(interfaceType)
				}
			}
		}
	}
	if err != nil {
		resp.Error = &proto.Error{
			Msg: err.Error(),
		}
	}
	return resp
}
