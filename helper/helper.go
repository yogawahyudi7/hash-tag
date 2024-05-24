package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yogawahyudi7/social-media/config"
	"golang.org/x/crypto/bcrypt"
)

func QueryIN(ids []int) (string, []any) {
	ar := []string{}
	er := []any{}
	for i, v := range ids {

		ar = append(ar, fmt.Sprintf("$%v", i+1))
		er = append(er, fmt.Sprintf("%v", v))
	}

	iQuery := strings.Join(ar, ",")

	return iQuery, er
}

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil

}

func ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(config *config.Server, id int, isAdmin bool) (string, error) {

	timeDuration, _ := time.ParseDuration(config.TokenDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"isAdmin": isAdmin,
		"exp":     config.TimeNow.Add(timeDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(config *config.Server, tokenString string) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
