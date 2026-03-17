package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/XRenso/barerepo/collector/internal/domain"
	"github.com/XRenso/barerepo/collector/internal/usecase"
	pb "github.com/XRenso/barerepo/pkg/pb"
)

type Handler struct {
	pb.UnimplementedRepoServiceServer
	uc *usecase.RepoUseCase
}

func NewHandler(uc *usecase.RepoUseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetRepo(ctx context.Context, req *pb.GetRepoRequest) (*pb.GetRepoResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return nil, status.Error(codes.InvalidArgument, "owner and repo are required")
	}

	repo, err := h.uc.GetRepo(ctx, req.Owner, req.Repo)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, domain.ErrForbidden):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetRepoResponse{
		Title:     repo.Title,
		Desc:      repo.Desc,
		Creator:   repo.Creator,
		StarsCnt:  int32(repo.StarsCnt),
		ForksCnt:  int32(repo.ForksCnt),
		IssuesCnt: int32(repo.IssuesCnt),
		Lang:      repo.Lang,
		Size:      int32(repo.Size),
		CreatedAt: repo.CreatedAt,
		License:   repo.License,
	}, nil
}
