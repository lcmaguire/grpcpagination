package grpcpagination

import (
	"context"
	"errors"
	"testing"

	"github.com/lcmaguire/grpcpagination/example"
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
		rpc := MockRpc{Res: res}
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
		rpc := MockRpc{Res: res}
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
		rpc := MockRpc{Err: expectedErr}
		err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool { return false })
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
	})

	t.Run("failure_rpc_error", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("err")
		req := &example.ListExamplesRequest{}
		rpc := MockRpc{Err: expectedErr}
		err := PaginateNextToken(ctx, req, rpc.ListExamples, func(ctx context.Context, response *example.ListExamplesResponse) bool { return false })
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
	})

}

type MockRpc struct {
	Res *example.ListExamplesResponse
	Err error
}

func (m *MockRpc) ListExamples(ctx context.Context, in *example.ListExamplesRequest, opts ...grpc.CallOption) (*example.ListExamplesResponse, error) {
	return m.Res, m.Err
}
