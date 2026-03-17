package collector

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/XRenso/barerepo/gateway/internal/domain"
	pb "github.com/XRenso/barerepo/pkg/pb"
)

type Client struct {
	grpc pb.RepoServiceClient
}

func NewClient(grpc pb.RepoServiceClient) *Client {
	return &Client{grpc: grpc}
}

func (c *Client) GetRepo(ctx context.Context, owner, repo string) (*domain.GHRepo, error) {
	resp, err := c.grpc.GetRepo(ctx, &pb.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("%w: %s", domain.ErrNotFound, st.Message())
			case codes.InvalidArgument:
				return nil, errors.New(st.Message())
			}
		}
		return nil, fmt.Errorf("collector: %w", err)
	}

	return &domain.GHRepo{
		Title:     resp.Title,
		Desc:      resp.Desc,
		Creator:   resp.Creator,
		StarsCnt:  resp.StarsCnt,
		ForksCnt:  resp.ForksCnt,
		IssuesCnt: resp.IssuesCnt,
		Lang:      resp.Lang,
		Size:      resp.Size,
		CreatedAt: resp.CreatedAt,
		License:   resp.License,
	}, nil
}
