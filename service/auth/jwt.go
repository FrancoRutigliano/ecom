package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/FrancoRutigliano/ecom/config"
	typePayload "github.com/FrancoRutigliano/ecom/types"
	"github.com/FrancoRutigliano/ecom/utils"
	"github.com/golang-jwt/jwt"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handlerFunc http.HandlerFunc, store typePayload.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the user request
		tokenString := getTokenFromRequest(r)
		// validate the jwt
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}
		// if is we need to fetch the userID from the db (id from the token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
		}

		// set context "  userID"
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		// call the function is th token is valid
		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteErr(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
