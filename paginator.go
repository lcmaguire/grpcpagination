// Package grpcpagination used to paginate a gRPC List endpoint that implement PageTokenRequest as it's request & NextPageTokenResponse as it's response.
package grpcpagination

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// GrpcFunc represents a List gRPC func that has a requestType and a response type.
type GrpcFunc[RequestType proto.Message, ResponseType proto.Message] func(ctx context.Context, request RequestType, opts ...grpc.CallOption) (ResponseType, error)

// FinishFunc used within PaginateNextToken to determine if pagination should be continued. Takes in the response from a List rpc call.
//
// to exit pagination return true from this func, to keep paginating return false.
type FinishFunc[ResponseType proto.Message] func(ctx context.Context, response ResponseType) bool

// PageTokenRequest represents a gRPC request where the page_token can be set. Requires protoc-gen-go-setters to be used.
type PageTokenRequest interface {
	proto.Message
	SetPageToken(in string)
}

// NextPageTokenResponse represents a gRPC response that will return the token for the next page.
type NextPageTokenResponse interface {
	proto.Message
	GetNextPageToken() string
}

// PaginateNextToken will paginate a gRPC endpoint so long as there is a next page, no exit condition met & no error returned.
func PaginateNextToken[RequestType PageTokenRequest, ResponseType NextPageTokenResponse](ctx context.Context, req RequestType, rpc GrpcFunc[RequestType, ResponseType], exitCondition FinishFunc[ResponseType]) error {
	for {
		res, err := rpc(ctx, req)
		if err != nil {
			return err
		}

		// if exit condition met
		if exitCondition(ctx, res) {
			break
		}

		// if no next page return.
		if res.GetNextPageToken() == "" {
			break
		}
		// set next page token
		req.SetPageToken(res.GetNextPageToken())
	}

	return nil
}
