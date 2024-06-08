package utils

import (
	"context"
	"errors"

	"github.com/SEC-Jobstreet/backend-job-service/domain/repository/model"
)

const (
	AuthorizationPayloadKey = "authorization_payload"
	CandidateRole           = "candidates"
	EmployerRole            = "employers"
	Admin                   = "admin"
)

func GetCurrentUser(ctx context.Context) (model.AuthClaim, error) {
	userCtx := ctx.Value(AuthorizationPayloadKey)
	currentUser, ok := userCtx.(model.AuthClaim)
	if !ok {
		return model.AuthClaim{}, errors.New("cannot parse to type User")
	}

	return currentUser, nil
}
