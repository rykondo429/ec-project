package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CognitoConfig Cognito設定
type CognitoConfig struct {
	UserPoolID string
	Region     string
	ClientID   string
	JWKSUrl    string // JSON Web Key Set URL
}

// CognitoMiddleware Cognito OpenID Connect認証ミドルウェア
func CognitoMiddleware(cfg CognitoConfig, logger *zap.Logger) echo.MiddlewareFunc {
	// JWKS (JSON Web Key Set) をキャッシュ
	jwksCache := NewJWKSCache(cfg.JWKSUrl, 1*time.Hour)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーから Bearer トークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// トークンなしの場合はアノニマスユーザーとして扱う
				c.Set("user_id", "anonymous")
				c.Set("email", "")
				return next(c)
			}

			// "Bearer " の部分を削除
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid authorization header format")
				return echo.NewHTTPError(http.StatusBadRequest, "invalid authorization header")
			}

			tokenString := parts[1]

			// トークンを検証
			claims, err := validateCognitoToken(c.Request().Context(), tokenString, cfg, jwksCache)
			if err != nil {
				logger.Warn("Token validation failed", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// クレームからユーザー情報を取得
			userID := claims["sub"].(string)
			email := ""
			if emailVal, ok := claims["email"].(string); ok {
				email = emailVal
			}

			// コンテキストに設定
			c.Set("user_id", userID)
			c.Set("email", email)
			c.Set("claims", claims)

			logger.Debug("Authenticated user",
				zap.String("user_id", userID),
				zap.String("email", email))

			return next(c)
		}
	}
}

// validateCognitoToken Cognitoトークンを検証
func validateCognitoToken(ctx context.Context, tokenString string, cfg CognitoConfig, jwksCache *JWKSCache) (jwt.MapClaims, error) {
	// トークンをパース (検証なし)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// kid (Key ID) を取得
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid in token header")
	}

	// JWKS から公開鍵を取得
	publicKey, err := jwksCache.GetKey(ctx, kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// トークンを検証
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// alg を確認
		if token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// クレームを検証
	if err := validateClaims(claims, cfg); err != nil {
		return nil, fmt.Errorf("claim validation failed: %w", err)
	}

	return claims, nil
}

// validateClaims クレームを検証
func validateClaims(claims jwt.MapClaims, cfg CognitoConfig) error {
	// iss (issuer) を検証
	expectedIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.Region, cfg.UserPoolID)
	if iss, ok := claims["iss"].(string); !ok || iss != expectedIssuer {
		return fmt.Errorf("invalid issuer: %v", claims["iss"])
	}

	// aud (audience) または client_id を検証
	if aud, ok := claims["aud"].(string); ok {
		if aud != cfg.ClientID {
			return fmt.Errorf("invalid audience: %s", aud)
		}
	} else if clientID, ok := claims["client_id"].(string); ok {
		if clientID != cfg.ClientID {
			return fmt.Errorf("invalid client_id: %s", clientID)
		}
	} else {
		return errors.New("missing aud or client_id")
	}

	// token_use を検証 (id または access)
	if tokenUse, ok := claims["token_use"].(string); ok {
		if tokenUse != "id" && tokenUse != "access" {
			return fmt.Errorf("invalid token_use: %s", tokenUse)
		}
	}

	// exp (expiration time) を検証
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return errors.New("token expired")
		}
	} else {
		return errors.New("missing exp")
	}

	return nil
}

// JWKS (JSON Web Key Set) キャッシュ
type JWKSCache struct {
	url       string
	cache     map[string]interface{}
	cacheTime time.Time
	ttl       time.Duration
}

// NewJWKSCache JWKSキャッシュを作成
func NewJWKSCache(url string, ttl time.Duration) *JWKSCache {
	return &JWKSCache{
		url:   url,
		cache: make(map[string]interface{}),
		ttl:   ttl,
	}
}

// GetKey kid から公開鍵を取得
func (j *JWKSCache) GetKey(ctx context.Context, kid string) (interface{}, error) {
	// キャッシュが有効かチェック
	if time.Since(j.cacheTime) > j.ttl || j.cache[kid] == nil {
		if err := j.refresh(ctx); err != nil {
			return nil, err
		}
	}

	key, ok := j.cache[kid]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", kid)
	}

	return key, nil
}

// refresh JWKS を取得して キャッシュを更新
func (j *JWKSCache) refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", j.url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	// 各キーをパース
	newCache := make(map[string]interface{})
	for _, key := range jwks.Keys {
		kid, ok := key["kid"].(string)
		if !ok {
			continue
		}

		// RSA公開鍵をパース
		publicKey, err := parseRSAPublicKey(key)
		if err != nil {
			continue
		}

		newCache[kid] = publicKey
	}

	j.cache = newCache
	j.cacheTime = time.Now()

	return nil
}

// parseRSAPublicKey JWK から RSA公開鍵をパース
func parseRSAPublicKey(jwk map[string]interface{}) (interface{}, error) {
	// kty が RSA であることを確認
	kty, ok := jwk["kty"].(string)
	if !ok || kty != "RSA" {
		return nil, errors.New("invalid key type")
	}

	// n (modulus) と e (exponent) を取得
	nStr, ok := jwk["n"].(string)
	if !ok {
		return nil, errors.New("missing n")
	}

	eStr, ok := jwk["e"].(string)
	if !ok {
		return nil, errors.New("missing e")
	}

	// Base64 URL デコード
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %w", err)
	}

	// big.Int に変換
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// RSA公開鍵を構築
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}
