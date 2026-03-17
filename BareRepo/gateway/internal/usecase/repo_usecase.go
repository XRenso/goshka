package usecase

import (
	"context"

	"github.com/XRenso/barerepo/gateway/internal/domain"
)

type RepoFetcher interface {
	GetRepo(ctx context.Context, owner, repo string) (*domain.GHRepo, error)
}

type RepoUseCase struct {
	fetcher RepoFetcher
}

func NewRepoUseCase(fetcher RepoFetcher) *RepoUseCase {
	return &RepoUseCase{fetcher: fetcher}
}

func (uc *RepoUseCase) GetRepo(ctx context.Context, owner, repo string) (*domain.GHRepo, error) {
	return uc.fetcher.GetRepo(ctx, owner, repo)
}
