package middleware

import (
	"context"
	"errors"
	"net/http"
	"root/lib/responser"
	"root/models"
	"root/services/login"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer = "avitoShop"
	prefix = "Bearer "
	ctxKey = "user"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.Contains(authHeader, prefix) {
			responser.SendHttpError(w, "Missing token", http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimPrefix(authHeader, prefix)

		u, err := verifyToken(rawToken)

		if err != nil {
			responser.SendHttpError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := putUserToContext(r.Context(), u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetUserFromContext(ctx context.Context) (models.User, bool) {
	u, ok := ctx.Value(ctxKey).(models.User)

	if !ok {
		return models.User{}, false
	}

	return u, true
}

func putUserToContext(ctx context.Context, u models.User) context.Context {
	return context.WithValue(ctx, ctxKey, u)
}

func verifyToken(rawtoken string) (models.User, error) {

	token, err := jwt.Parse(
		rawtoken,
		keyFunc(),
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		return models.User{}, err
	}

	if !token.Valid {
		return models.User{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return models.User{}, errors.New("invalid token")
	}

	uId := int(claims["user_id"].(float64))
	uName := claims["user_name"].(string)

	return models.User{
		ID:       uId,
		Username: uName,
	}, nil
}

func keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (any, error) { return []byte(login.Secret), nil }
}
