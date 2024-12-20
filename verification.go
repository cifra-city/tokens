package tokens

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// VerifyJWTAndExtractClaims validates a JWT token and extracts relevant claims.
func VerifyJWTAndExtractClaims(
	ctx context.Context,
	tokenString,
	secretKey string,
	log *logrus.Logger) (
	userId uuid.UUID,
	deviceId uuid.UUID,
	tokenVersion int,
	role string,
	err error,
) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		log.Debugf("Token parsing failed: %v", err)
		return uuid.Nil, uuid.Nil, 0, "", err
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Debugf("Invalid user ID in claims: %v", err)
		return uuid.Nil, uuid.Nil, 0, "", err
	}

	// Extract token version and role
	tokenVersion = claims.TokenVersion
	role = claims.Role
	deviceId = claims.DeviceID

	if tokenVersion == 0 {
		log.Debug("Token version is missing in claims")
		return uuid.Nil, uuid.Nil, 0, "", jwt.ErrTokenMalformed
	}

	return userID, deviceId, tokenVersion, role, nil
}
