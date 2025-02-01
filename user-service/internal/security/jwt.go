package security

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"user-service/proto/pb"
)

var jwtKey = []byte("qp_TGOFe56TIehvKUOzQAuMVEqelvKgWR9sznKmPrxBLRLZfdgsngdgzEIfdyQuzQeMhysnScNVBB5qwAuPbt29_IUbEx1V5r5eybrbkoDJdLpvQFUubvzULjqZUTKmlZ")

type Claims struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	CompanyId int32  `json:"company_id"`
	jwt.StandardClaims
}

func GenerateToken(user *pb.GetUserByIdResponse) (string, error) {
	expirationTime := time.Now().Add(24 * 10 * time.Hour)
	claims := &Claims{
		Username:  user.Id,
		Role:      user.Role,
		CompanyId: user.CompanyId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
