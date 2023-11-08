package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CreateJWT takes a user ID and returns a JWT token string
// The token contains an expiresAt and userId claims
func CreateJWT(userId int) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Create the Claims
	claims := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"userId":    userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

func GetClaimFromJWT(tokenString string, claimName string) (interface{}, error) {
	token, err := ValidateJWT(tokenString) // Implement the ValidateJWT function to parse and validate the token.
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to extract claims")
	}

	claimValue, exists := claims[claimName]
	if !exists {
		return nil, fmt.Errorf("Claim not found")
	}

	return claimValue, nil
}

func GetUserIdFromJWT(tokenString string) (int, error) {
	claimValue, err := GetClaimFromJWT(tokenString, "userId")
	if err != nil {
		return 0, err
	}

	userId, ok := claimValue.(float64)
	if !ok {
		return 0, fmt.Errorf("userId claim is not a valid number")
	}

	return int(userId), nil
}
