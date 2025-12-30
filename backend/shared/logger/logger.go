package logger

import "go.uber.org/zap"

// NewLogger ロガーを作成
func NewLogger(environment string) (*zap.Logger, error) {
	if environment == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

// NewSugaredLogger シュガーロガーを作成
func NewSugaredLogger(environment string) (*zap.SugaredLogger, error) {
	logger, err := NewLogger(environment)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
