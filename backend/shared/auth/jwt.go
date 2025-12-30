package auth

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTClaims JWT クレーム
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// JWTMiddleware JWT検証ミドルウェア
func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーから Bearer トークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// トークンなしの場合はアノニマスユーザーとして扱う
				c.Set("user_id", "anonymous")
				return next(c)
			}

			// "Bearer " の部分を削除
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(400, "invalid authorization header")
			}

			tokenString := parts[1]

			// トークンをパース
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return echo.NewHTTPError(401, "invalid token")
			}

			if !token.Valid {
				return echo.NewHTTPError(401, "invalid token")
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return echo.NewHTTPError(401, "invalid claims")
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// GenerateToken トークン生成 (テスト用)
func GenerateToken(userID, email, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: userID,
		Email:  email,
	})

	return token.SignedString([]byte(secretKey))
}
