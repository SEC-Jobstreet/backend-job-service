package utils

import (
	"context"
	"errors"

	"github.com/SEC-Jobstreet/backend-job-service/models"
)

const (
	AuthorizationPayloadKey = "authorization_payload"
	CandidateRole           = "candidates"
	EmployerRole            = "employers"
	Admin                   = "admin"
)

func GetCurrentUser(ctx context.Context) (models.AuthClaim, error) {
	userCtx := ctx.Value(AuthorizationPayloadKey)
	currentUser, ok := userCtx.(models.AuthClaim)
	if !ok {
		return models.AuthClaim{}, errors.New("cannot parse to type User")
	}

	return currentUser, nil
}
