package server

import (
	"github.com/haileyok/cocoon/internal/helpers"
	"github.com/haileyok/cocoon/models"
	"github.com/labstack/echo/v4"
)

type ComAtprotoRepoApplyWritesRequest struct {
	Repo       string                          `json:"repo" validate:"required,atproto-did"`
	Validate   *bool                           `json:"bool,omitempty"`
	Writes     []ComAtprotoRepoApplyWritesItem `json:"writes"`
	SwapCommit *string                         `json:"swapCommit"`
}

type ComAtprotoRepoApplyWritesItem struct {
	Type       string          `json:"$type"`
	Collection string          `json:"collection"`
	Rkey       string          `json:"rkey"`
	Value      *MarshalableMap `json:"value,omitempty"`
}

func (s *Server) handleApplyWrites(e echo.Context) error {
	repo := e.Get("repo").(*models.RepoActor)

	var req ComAtprotoRepoApplyWritesRequest
	if err := e.Bind(&req); err != nil {
		s.logger.Error("error binding", "error", err)
		return helpers.ServerError(e, nil)
	}

	if err := e.Validate(req); err != nil {
		s.logger.Error("error validating", "error", err)
		return helpers.InputError(e, nil)
	}

	if repo.Repo.Did != req.Repo {
		s.logger.Warn("mismatched repo/auth")
		return helpers.InputError(e, nil)
	}

	ops := []Op{}
	for _, item := range req.Writes {
		ops = append(ops, Op{
			Type:       OpType(item.Type),
			Collection: item.Collection,
			Rkey:       &item.Rkey,
			Record:     item.Value,
		})
	}

	if err := s.repoman.applyWrites(repo.Repo, ops, req.SwapCommit); err != nil {
		s.logger.Error("error applying writes", "error", err)
		return helpers.ServerError(e, nil)
	}

	return nil
}
