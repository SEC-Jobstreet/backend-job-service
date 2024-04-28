package gapi

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

type Auth struct {
	jwk               *JWK
	jwkURL            string
	cognitoRegion     string
	cognitoUserPoolID string
}

type Config struct {
	CognitoRegion     string
	CognitoUserPoolID string
}

type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

func NewAuth(config utils.Config) *Auth {
	a := &Auth{
		cognitoRegion:     config.CognitoRegionEmployers,
		cognitoUserPoolID: config.CognitoUserPoolIDEmployers,
	}

	a.jwkURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.cognitoRegion, a.cognitoUserPoolID)
	err := a.CacheJWK()
	if err != nil {
		logrus.Fatal(err)
	}

	return a
}

func (m *Auth) CacheJWK() error {
	req, err := http.NewRequest("GET", m.jwkURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jwk := new(JWK)
	err = json.Unmarshal(body, jwk)
	if err != nil {
		return err
	}

	m.jwk = jwk
	return nil
}

func (m *Auth) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key := convertKey(m.jwk.Keys[1].E, m.jwk.Keys[1].N)
		return key, nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}

func (m *Auth) JWK() *JWK {
	return m.jwk
}

func convertKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}

func (server *Server) authorizeUser(ctx context.Context, config utils.Config, accessibleRoles []string) (*models.AuthClaim, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]

	// Validate token cognito
	auth := NewAuth(config)
	err := auth.CacheJWK()
	if err != nil {
		return nil, fmt.Errorf("AuthMiddleware - Error cacheJWK, error = %v", err)
	}

	token, err := auth.ParseJWT(accessToken)
	if err != nil {
		return nil, fmt.Errorf("AuthMiddleware - Error ParseJWT, error = %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("AuthMiddleware - Invalid token")
	}

	// Parse token to struct user and store into context
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("AuthMiddleware - Invalid token")
	}
	var currentUser models.AuthClaim
	errDecode := mapstructure.Decode(claims, &currentUser)
	if errDecode != nil {
		return nil, fmt.Errorf("AuthMiddleware - Cannot decode claims %v to current user, error %v", utils.LogFull(claims), errDecode)
	}

	if len(accessibleRoles) != 0 && !HasPermission(currentUser.CognitoGroups, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return &currentUser, nil
}

func HasPermission(userRoles []string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		for _, userRole := range userRoles {
			if userRole == role {
				return true
			}
		}
	}
	return false
}
