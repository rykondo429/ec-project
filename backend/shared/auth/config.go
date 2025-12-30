package auth

import (
	"fmt"
	"os"
)

// GetCognitoConfig 環境変数からCognito設定を取得
func GetCognitoConfig() CognitoConfig {
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	region := getEnvOrDefault("COGNITO_REGION", "us-east-1")
	clientID := os.Getenv("COGNITO_CLIENT_ID")

	// JWKS URL を構築
	jwksURL := os.Getenv("COGNITO_JWKS_URL")
	if jwksURL == "" {
		// LocalStackの場合
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			jwksURL = fmt.Sprintf("%s/", endpoint)
		} else {
			// 本番環境
			jwksURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
		}
	}

	return CognitoConfig{
		UserPoolID: userPoolID,
		Region:     region,
		ClientID:   clientID,
		JWKSUrl:    jwksURL,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
