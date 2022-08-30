package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"neft.web/models"
)

var jwtKey = []byte("Sabadell0310JED")

type JWTClaim struct {
	RemmemberHash string `json:"remmemberHash"`
	RoleID        int    `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(remmemberHash string, roleID int) (tokenString string, err error) {
	expirationTime := time.Now().Local().Add(time.Duration(12 * time.Hour))
	claims := &JWTClaim{
		RemmemberHash: remmemberHash,
		RoleID:        roleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		fmt.Println("Hola")
		err = models.ERR_JWT_CLAIMS_INVALID
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		fmt.Println("Adios")
		err = models.ERR_JWT_TOKEN_EXPIRED
		return
	}
	return
}

func ReturnClaims(signedToken string) (claim *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = models.ERR_JWT_CLAIMS_INVALID
		return nil, err
	}

	return claims, nil
}
