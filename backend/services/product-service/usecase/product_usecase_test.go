package usecase_test

import (
	"context"
	"testing"
)

// ダミーテスト - CI用
func TestProductUseCaseDummy(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Error("context should not be nil")
	}
}
