package grpcpagination

import (
	"context"
	"errors"
	"testing"

	"github.com/lcmaguire/protoc-gen-go-setters/example"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestPaginateNextToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success_one_page", func(t *testing.T) {
		t.Parallel()
		req := &example.ListExamplesRequest{}
		res := &example.ListExamplesResponse{}
		rpc := mockRpc{res: res}
		count := 0
		expectedCount := 1
		err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool {
			count++
			return count == expectedCount
		})
		require.NoError(t, err)
		require.Equal(t, expectedCount, count)
	})

	t.Run("success_five_pages", func(t *testing.T) {
		t.Parallel()
		req := &example.ListExamplesRequest{}
		expectedToken := "page_token"
		res := &example.ListExamplesResponse{
			NextPageToken: expectedToken,
		}
		rpc := mockRpc{res: res}
		count := 0
		expectedCount := 5
		err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool {
			count++
			return count == expectedCount
		})
		require.NoError(t, err)
		require.Equal(t, expectedCount, count)
		require.Equal(t, expectedToken, req.PageToken)
	})

	t.Run("failure_rpc_error", func(t *testing.T) {
		t.Parallel()
		req := &example.ListExamplesRequest{}
		expectedErr := errors.New("err")
		rpc := mockRpc{err: expectedErr}
		err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool { return false })
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
	})

	t.Run("failure_rpc_error", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("err")
		req := &example.ListExamplesRequest{}
		rpc := mockRpc{err: expectedErr}
		err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool { return false })
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
	})

}

type mockRpc struct {
	res *example.ListExamplesResponse
	err error
}

func (m *mockRpc) ListExamples(ctx context.Context, in *example.ListExamplesRequest, opts ...grpc.CallOption) (*example.ListExamplesResponse, error) {
	return m.res, m.err
}
