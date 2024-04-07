package grpcpagination

import (
	"context"
	"fmt"

	"github.com/lcmaguire/protoc-gen-go-setters/example"
)

func ExamplePaginateNextToken() {
	// use your own rpcClient
	expectedToken := "page_token"
	res := &example.ListExamplesResponse{
		NextPageToken: expectedToken,
	}
	rpc := mockRpc{res: res}
	ctx := context.Background()

	req := &example.ListExamplesRequest{}

	const maxPages = 5
	pages := 0
	err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool {
		pages++
		return pages == maxPages
	})

	fmt.Println("err will only exist if there is a rpc error", err)
	fmt.Println("pages visited", pages)

	// Output:
	// err will only exist if there is a rpc error <nil>
	// pages visited 5
}
