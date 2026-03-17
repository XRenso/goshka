package usecase

import (
	"context"

	"github.com/XRenso/barerepo/collector/internal/domain"
)

type RepoProvider interface {
	FetchRepo(ctx context.Context, owner, repo string) (*domain.GHRepo, error)
}

type RepoUseCase struct {
	provider RepoProvider
}

func NewRepoUseCase(provider RepoProvider) *RepoUseCase {
	return &RepoUseCase{provider: provider}
}

func (uc *RepoUseCase) GetRepo(ctx context.Context, owner, repo string) (*domain.GHRepo, error) {
	return uc.provider.FetchRepo(ctx, owner, repo)
}
