package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/XRenso/barerepo/gateway/internal/domain"
	"github.com/XRenso/barerepo/gateway/internal/usecase"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	uc *usecase.RepoUseCase
}

func NewHandler(uc *usecase.RepoUseCase) *Handler {
	return &Handler{uc: uc}
}

// GetRepo godoc
// @Summary      Get repository info
// @Description  Returns info about a GitHub repository (title, description, stars, forks, issues, language, size, license, created_at)
// @Tags         repositories
// @Produce      json
// @Param        owner  path      string  true  "Repository owner"
// @Param        repo   path      string  true  "Repository title"
// @Success      200    {object}  domain.GHRepo
// @Failure      400    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /api/v1/repos/{owner}/{repo} [get]
func (h *Handler) GetRepo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "owner and repo are required"})
		return
	}

	result, err := h.uc.GetRepo(c.Request.Context(), owner, repo)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
